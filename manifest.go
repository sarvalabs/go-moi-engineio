package engineio

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sarvalabs/go-polo"
	"golang.org/x/crypto/blake2b"
	"gopkg.in/yaml.v3"
)

// Encoding is an enum with variants that describe
// encoding schemes supported for file objects.
type Encoding int

const (
	POLO Encoding = iota
	JSON
	YAML
)

// Manifest is the canonical deployment artifact for logics in MOI.
//
// It is a composite artifact that describes the bytecode, the binary interface (ABI) and
// other parameters for the runtime of choice. The spec for LogicManifests is available at
// https://sarvalabs.notion.site/Logic-Manifest-Standard-93f5fee1af8d4c3cad155b9827b97930?pvs=4
type Manifest struct {
	Syntax   string            `yaml:"syntax" json:"syntax"`
	Engine   ManifestEngine    `yaml:"engine" json:"engine"`
	Elements []ManifestElement `yaml:"elements" json:"elements"`
}

// ManifestEngine describes the engine specific information in the Manifest
type ManifestEngine struct {
	Kind  string   `yaml:"kind" json:"kind"`
	Flags []string `yaml:"flags" json:"flags"`
}

// ManifestElement describes a single element in the Manifest.
// It is converted into a LogicElement after compilation.
//
// Each element is of a particular type (described by the engine runtime) and is identified by
// a unique 64-bit pointer  and describes its dependencies with other elements in the manifest.
//
// The data of the manifest element must be resolved into the format specific for the runtime based on its
// kind. The raw object to decode into can be accessed with the GetElementGenerator method of EngineRuntime.
type ManifestElement struct {
	Ptr  ElementPtr            `yaml:"ptr" json:"ptr"`
	Deps []ElementPtr          `yaml:"deps" json:"deps"`
	Kind ElementKind           `yaml:"kind" json:"kind"`
	Data ManifestElementObject `yaml:"data" json:"data"`
}

// The ManifestElementObject is a placeholder that we can decode the element's data into.
// Types that implement them are specified by the runtime and must de/serializable to all supported formats.
type ManifestElementObject interface {
	polo.Polorizable
	polo.Depolorizable
}

// ManifestElementGenerator is a generator function that returns an empty instance of the element type
type ManifestElementGenerator func() ManifestElementObject

// NewManifest decodes the given raw data of the specified encoding type into a Manifest.
// Fails if the encoding is unsupported or if the data is malformed.
func NewManifest(data []byte, encoding Encoding) (*Manifest, error) {
	manifest := new(Manifest)

	switch encoding {
	case JSON:
		if err := json.Unmarshal(data, manifest); err != nil {
			return nil, err
		}
	case POLO:
		if err := polo.Depolorize(manifest, data); err != nil {
			return nil, err
		}
	case YAML:
		if err := yaml.Unmarshal(data, manifest); err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("unsupported manifest encoding")
	}

	return manifest, nil
}

// ReadManifestFile reads a file at the specified filepath and decodes it into a Manifest.
// The encoding format of the file is determined from the file extension.
func ReadManifestFile(path string) (*Manifest, error) {
	path, _ = filepath.Abs(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.Errorf("manifest file not found @ '%v'", path)
	}

	var (
		extension string
		encoding  Encoding
	)

	switch extension = filepath.Ext(path); extension {
	case ".json":
		encoding = JSON
	case ".polo":
		encoding = POLO
	case ".yaml":
		encoding = YAML
	default:
		return nil, errors.Errorf("manifest file has unsupported extension: '%v'", extension)
	}

	encoded, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read manifest file")
	}

	manifest, err := NewManifest(encoded, encoding)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode %v manifest data", extension)
	}

	return manifest, nil
}

// Hash returns the 256-bit hash of the Manifest.
// The hash is derived by applying the Blake2b hashing
// function on the POLO encoded bytes of the Manifest.
func (manifest Manifest) Hash() ([32]byte, error) {
	encoded, err := manifest.Encode(POLO)
	if err != nil {
		return [32]byte{}, err
	}

	return blake2b.Sum256(encoded), nil
}

// Encode returns the encoded bytes form of the Manifest for the specified encoding.
func (manifest Manifest) Encode(encoding Encoding) ([]byte, error) {
	switch encoding {
	case JSON:
		return json.Marshal(manifest)
	case POLO:
		return polo.Polorize(manifest)
	case YAML:
		return yaml.Marshal(manifest)

	default:
		return nil, errors.New("unsupported manifest encoding")
	}
}

// Header returns the header information of the Manifest as a ManifestHeader
func (manifest Manifest) Header() ManifestHeader {
	return ManifestHeader{manifest.Syntax, manifest.Engine}
}

// ManifestHeader represents the header for a Manifest and describes its syntax form
// and engine specification. Useful for determining which engine to use to handle the
// Manifest. Every engine's manifest implementation must be able to decode into this header.
type ManifestHeader struct {
	Syntax string         `yaml:"syntax" json:"syntax"`
	Engine ManifestEngine `yaml:"engine" json:"engine"`
}

// LogicEngine returns the normalized form of the logic engine value in the ManifestHeader.
// It is capitalized to uppercase letter and converted into an EngineKind
func (header ManifestHeader) LogicEngine() EngineKind {
	return EngineKind(strings.ToUpper(header.Engine.Kind))
}

func (header ManifestHeader) validate() error {
	if header.Syntax != "0.1.0" {
		return errors.New("unsupported manifest syntax")
	}

	if _, ok := FetchRuntime(header.LogicEngine()); !ok {
		return errors.New("unsupported manifest engine: element registry not found")
	}

	return nil
}

func (manifest *Manifest) Depolorize(depolorizer *polo.Depolorizer) (err error) {
	type ManifestPOLO struct {
		Syntax   string
		Engine   ManifestEngine
		Elements []struct {
			Ptr  ElementPtr
			Deps []ElementPtr
			Kind ElementKind
			Data polo.Any
		}
	}

	raw := new(ManifestPOLO)
	if err = depolorizer.Depolorize(raw); err != nil {
		return err
	}

	manifest.Syntax = raw.Syntax
	manifest.Engine = raw.Engine

	if err = manifest.Header().validate(); err != nil {
		return err
	}

	runtime, _ := FetchRuntime(manifest.Header().LogicEngine())

	manifest.Elements = make([]ManifestElement, 0, len(raw.Elements))

	for _, element := range raw.Elements {
		generator, ok := runtime.GetElementGenerator(element.Kind)
		if !ok {
			return errors.Errorf("unrecognized element kind: '%v'", element.Kind)
		}

		elementDepolorizer, err := polo.NewDepolorizer(element.Data)
		if err != nil {
			return err
		}

		object := generator()
		if err = object.Depolorize(elementDepolorizer); err != nil {
			return err
		}

		manifest.Elements = append(manifest.Elements, ManifestElement{
			Ptr:  element.Ptr,
			Kind: element.Kind,
			Deps: element.Deps,
			Data: object,
		})
	}

	return nil
}

func (manifest *Manifest) UnmarshalJSON(data []byte) (err error) {
	type ManifestJSON struct {
		Syntax   string         `json:"syntax"`
		Engine   ManifestEngine `json:"engine"`
		Elements []struct {
			Ptr  ElementPtr      `json:"ptr"`
			Deps []ElementPtr    `json:"deps"`
			Kind ElementKind     `json:"kind"`
			Data json.RawMessage `json:"data"`
		} `json:"elements"`
	}

	raw := new(ManifestJSON)
	if err = json.Unmarshal(data, raw); err != nil {
		return err
	}

	manifest.Syntax = raw.Syntax
	manifest.Engine = raw.Engine

	if err = manifest.Header().validate(); err != nil {
		return err
	}

	runtime, _ := FetchRuntime(manifest.Header().LogicEngine())

	manifest.Elements = make([]ManifestElement, 0, len(raw.Elements))

	for _, element := range raw.Elements {
		generator, ok := runtime.GetElementGenerator(element.Kind)
		if !ok {
			return errors.Errorf("unrecognized element kind: '%v'", element.Kind)
		}

		object := generator()
		if err = json.Unmarshal(element.Data, object); err != nil {
			return err
		}

		manifest.Elements = append(manifest.Elements, ManifestElement{
			Ptr:  element.Ptr,
			Kind: element.Kind,
			Deps: element.Deps,
			Data: object,
		})
	}

	return nil
}

func (manifest *Manifest) UnmarshalYAML(node *yaml.Node) error {
	type ManifestYAML struct {
		Syntax   string         `yaml:"syntax"`
		Engine   ManifestEngine `yaml:"engine"`
		Elements []struct {
			Ptr  ElementPtr   `yaml:"ptr"`
			Deps []ElementPtr `yaml:"deps"`
			Kind ElementKind  `yaml:"kind"`
			Data yaml.Node    `yaml:"data"`
		} `yaml:"elements"`
	}

	raw := new(ManifestYAML)
	if err := node.Decode(raw); err != nil {
		return err
	}

	manifest.Syntax = raw.Syntax
	manifest.Engine = raw.Engine

	if err := manifest.Header().validate(); err != nil {
		return err
	}

	runtime, _ := FetchRuntime(manifest.Header().LogicEngine())

	manifest.Elements = make([]ManifestElement, 0, len(raw.Elements))

	for _, element := range raw.Elements {
		generator, ok := runtime.GetElementGenerator(element.Kind)
		if !ok {
			return errors.Errorf("unrecognized element kind: '%v'", element.Kind)
		}

		object := generator()
		if err := element.Data.Decode(object); err != nil {
			return err
		}

		manifest.Elements = append(manifest.Elements, ManifestElement{
			Ptr:  element.Ptr,
			Kind: element.Kind,
			Deps: element.Deps,
			Data: object,
		})
	}

	return nil
}

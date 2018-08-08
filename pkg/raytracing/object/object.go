package object

import (
	"encoding/json"
	"fmt"

	"github.com/BrendanBurkhart/raytracer/pkg/raytracing"
)

// Object provides an interface for intersecting with 3D objects and their materials
type Object interface {
	Intersect(r raytracing.Ray, maxRange float64) (bool, float64)
	SurfaceNormal(point raytracing.Vector) raytracing.Vector
	MaterialID() int
}

// Material can be embedded in an object so it satisfies the MaterialID getter requirement of Object
type Material struct {
	Material int
}

// MaterialID returns the id of the material attached to the object
func (om *Material) MaterialID() int {
	return om.Material
}

// shapeUnmarshaller unmarshals JSON data into a specific implementation of Object
type objectFactory func(*json.RawMessage) (Object, error)

var objectFactoryMap = map[string]objectFactory{
	"plane":  planeFactory,
	"sphere": sphereFactory,
	"box":    boxFactory,
}

// JSONObjects is a named type to allow a slice of interfaces to have custom JSON unmarshalling
type JSONObjects []Object

// UnmarshalJSON allows an array of different structs which all implement Object to be unmarhsal to an array of Object
func (jsonObjects *JSONObjects) UnmarshalJSON(b []byte) error {
	var rawObjects []*json.RawMessage
	if err := json.Unmarshal(b, &rawObjects); err != nil {
		return err
	}

	var typingData []map[string]*json.RawMessage
	if err := json.Unmarshal(b, &typingData); err != nil {
		return err
	}

	for i, typing := range typingData {
		obj, err := unmarshalObject(typing, rawObjects[i])
		if err != nil {
			return err
		}
		*jsonObjects = append(*jsonObjects, obj)
	}
	return nil
}

// findObjectFactory returns the correct factory function for shape primitive (Object) based on its typing data
func findObjectFactory(typing map[string]*json.RawMessage) (factory objectFactory, err error) {
	rawShapeType, ok := typing["type"]
	if !ok {
		return nil, fmt.Errorf("JSON object does not contain key 'type' needed to unmarshal it")
	}

	var shapeType string
	if err = json.Unmarshal(*rawShapeType, &shapeType); err != nil {
		return nil, fmt.Errorf("error unmarshalling shape type to string: %v", err)
	}

	factory, ok = objectFactoryMap[shapeType]
	if !ok {
		return nil, fmt.Errorf("cannot find object type %s referenced in JSON data", shapeType)
	}
	return factory, nil
}

// unmarshalObject unmarshals the data into a struct implementing Object and returns it as an Object
func unmarshalObject(typing map[string]*json.RawMessage, data *json.RawMessage) (Object, error) {
	factory, err := findObjectFactory(typing)
	if err != nil {
		return nil, err
	}

	var obj Object
	obj, err = factory(data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

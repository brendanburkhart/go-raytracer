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

// Custom JSON unmarshalling to handle different implements of the Object interface

// JSONObjects is a named type to allow a slice of interfaces to have custom JSON unmarshalling
type JSONObjects []Object

// UnmarshalJSON allows an array of different structs which all implement Object to be unmarhsal to an array of Object
func (jsonObjects *JSONObjects) UnmarshalJSON(b []byte) error {
	var rawObjects []*json.RawMessage
	if err := json.Unmarshal(b, &rawObjects); err != nil {
		return err
	}

	var types []map[string]*json.RawMessage
	if err := json.Unmarshal(b, &types); err != nil {
		return err
	}

	*jsonObjects = append(*jsonObjects)

	for i, typing := range types {
		obj, err := unmarshalObject(typing, rawObjects[i])
		if err != nil {
			return err
		}
		*jsonObjects = append(*jsonObjects, obj)
	}
	return nil
}

// unmarshalObject unmarshals the data into a struct implementing Object and returns it as an Object
func unmarshalObject(typing map[string]*json.RawMessage, data *json.RawMessage) (Object, error) {
	for key, value := range typing {
		if key == "type" {
			var shapeType string
			if err := json.Unmarshal(*value, &shapeType); err != nil {
				return nil, err
			}
			switch shapeType {
			case "sphere":
				obj := Sphere{}
				if err := json.Unmarshal(*data, &obj); err != nil {
					return nil, err
				}
				return obj, nil
			case "plane":
				obj := Plane{}
				if err := json.Unmarshal(*data, &obj); err != nil {
					return nil, err
				}
				obj.Normalize()
				return obj, nil
			case "box":
				obj := Box{}
				if err := json.Unmarshal(*data, &obj); err != nil {
					return nil, err
				}
				obj.Initialize()
				return obj, nil
			default:
				return nil, fmt.Errorf("cannot find object type %s referenced in json data", shapeType)
			}
		}
	}
	return nil, fmt.Errorf("JSON object does not contain key 'type' needed to unmarshal it")
}

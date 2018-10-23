# Raytracer

This Raytracer is written in Go, and renders .json files into an image. I wrote this to learn more about raytracing, and to try out Go on something other than server-side programming.

### Features

This renders large/complicated scenes quickly using goroutines. Supports orthographic, simple perspective and fisheye projections.

It currently only supports materials, and does not support rendering textures or UV mapping. Additionally, only planes, spheres and boxes are supported. I may add support for UV mapping and more complex/custom shapes eventually.

Both Lambertian and Phong lighting models are supported, and both work with reflections and shadows. Refraction is not currently available.

I added simple configurable anti-aliasing through super sampling.

### Installing

Requires Go to be installed and setup. If it is not, refer to golang.org/doc/install to get started.

Run `go get -u github.com/BrendanBurkhart/raytracer/...` to install or update.

### Usage

Run the executable with the data file(s) and/or folder(s) containing the scenes to be rendered. Each scene will be rendered and output into a PNG of the same name as the scene's data file in the same location. Example: `raytracing.exe ./scenes ./example.json`.

## Scene data description

Scenes are described using JSON files in the following format:

```{
  "width": Output image width,
  "height": Output image height,
  "camera": {
    "position": Vector, specifies camera origin,
    "target": Vector, specifies where the camera is pointed,
    "roll": Camera roll in degrees, positive is counter-clockwise when facing the same direction as the camera,

    "antiAliasingFactor": Super samples per pixel, must be at least 1. Optional, default is 1,
    "lightingModel": One of "lambertian", or "phong". Optional, default is "phong",

    "projection": Projection type - one of "perspective", "orthographic", or "fisheye",

    "viewWidth": If using an orthographic projection, viewWidth must be specified. It is the view width of the rendered image in in-scene units. Can be used with a perspective projection, in which case focalLength must be specified.
    "hfov": If using a fisheye projection, hfov must be specified. It is the horizontal field of view in degrees. Can optionally replace viewWidth for a perspective projection.
    "focalLength": For perspective projection, distances from origin to render-plane. If this is not specified, opticalRadius must be.
    "opticalRadius": Radius of circle around camera origin in which render-plane is fit as plane with angle matching hfov.
  },
  "scene": {
    "materials": [Materials],
    "lights": [Lights],
    "objects": [Object primitives]
  }
}
```

Vectors are specified as `{"x": x, "y": y, "z": z}`.
Colors are specified as `{"red": 0.0-1.0, "green": 0.0-1.0, "blue": 0.0-1.0}`.

Materials are specified as:

```JSON
{
    "specular": Specular color,
    "diffuse": Diffuse color,
    "ambient": Ambient color,
    "alpha": 0 or greater, higher values create brighter, smaller specular highlights,
    "reflectance": 0.0 or greater, percentage of light reflected by material
},
```

Lights are specified as:

```JSON
{
    "position": Vector,
    "specular": Specular component, color,
    "diffuse": Diffuse component, color,
    "ambient": Ambient component, color
}
```

Object in the scene can be one of three primitives: sphere, box or plane.

Sphere:

```JSON
{
    "type": "sphere",
    "center": Position vector,
    "radius": Radius of sphere,
    "material": Index of material within array of materials
}
```

Box:

```JSON
{
    "type": "box",
    "minCorner": Position vector of minimum corner,
    "maxCorner": Position vector of maximum corner,
    "material": Index of material within array of materials
},
```

Plane:

```JSON
{
    "type": "plane",
    "point": Position vector of any point in plane,
    "normal": Normal vector of plane,
    "material": Index of material within array of materials
},
```

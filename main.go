package main

import (
	"fmt"
	"math"
	"github.com/gen2brain/raylib-go/raylib"
)

type light_source struct {
	light string
	intensity float32
	position rl.Vector3
	direction rl.Vector3
}

func new_light_source(light string, intensity float32) *light_source {
	return &light_source{
		light: light,
		intensity: intensity,
		position: rl.NewVector3(0, 0, 0),
		direction: rl.NewVector3(0, 0, 0),
	}
}

type sphere struct {
	center rl.Vector3
	radius float32
	color rl.Vector3
}

func color_transform (col rl.Vector3) rl.Color {
	return rl.ColorFromNormalized(rl.NewVector4(col.X, col.Y, col.Z, 1))
}

func canvas_to_viewport(u int32, v int32, ratio float32) rl.Vector3 {
	return rl.NewVector3(float32(u) * ratio, float32(v) * ratio, VIEWPORT_DIST)
}

func intersect_sphere(origin rl.Vector3, ray rl.Vector3, sph sphere) (float32, float32) {
	rad := sph.radius
	OC := rl.Vector3Subtract(origin, sph.center)

	a := rl.Vector3DotProduct(ray, ray)
	b := 2 * rl.Vector3DotProduct(OC, ray)
	c := rl.Vector3DotProduct(OC, OC) - rad * rad

	discriminant := b * b - 4 * a * c
	if (discriminant < 0) {
		return r_max, r_max
	} else {
		sqrt_dis := float32(math.Sqrt(float64(discriminant)))
		r1 := (-b + sqrt_dis) / 2 / a
		r2 := (-b - sqrt_dis) / 2 / a
		return r1, r2
	}
}

func trace_ray(origin rl.Vector3, ray rl.Vector3, scene []sphere, lights []light_source) rl.Vector3 {
	closest_r := r_max
	color := rl.NewVector3(0, 0, 0)
	sphere_num := len(scene)

	for s := 0; s < sphere_num; s++ {
		r1, r2 := intersect_sphere(origin, ray, scene[s])
		if (r1 > r_min && r1 < r_max && r1 < closest_r) {
			closest_r = r1
			point := rl.Vector3Add(origin, rl.Vector3Multiply(ray, closest_r))
			intensity := let_there_be_light(point, scene[s], lights)
			color = rl.Vector3Multiply(scene[s].color , intensity)
		}
		if (r2 > r_min && r2 < r_max && r2 < closest_r) {
			closest_r = r2
			point := rl.Vector3Add(origin, rl.Vector3Multiply(ray, closest_r))
			intensity := let_there_be_light(point, scene[s], lights)
			color = rl.Vector3Multiply(scene[s].color , intensity)
		}

	}

	return color
}

func let_there_be_light(
	point rl.Vector3, 
	sph sphere, 
	lights []light_source) float32 {
	var intensity float32 = 0

	normal := rl.Vector3Normalize(rl.Vector3Subtract(point, sph.center))
	ldir := normal
	var n_dot_l float32 = 1

	for l := 0; l < len(lights); l++ {
		if lights[l].light == "ambient" {
			intensity += lights[l].intensity
		} else {
			if lights[l].light == "point" {
				ldir = rl.Vector3Normalize(rl.Vector3Subtract(lights[l].position, point)) 
				n_dot_l = rl.Vector3DotProduct(normal, ldir)
				if n_dot_l > 0 {
					intensity += lights[l].intensity * n_dot_l
				}
			} 
			if lights[l].light == "directional" {
				ldir = lights[l].direction
				n_dot_l = rl.Vector3DotProduct(normal, ldir)
				if n_dot_l > 0 {
					intensity += lights[l].intensity * n_dot_l
				}
			}
		}
		
	}

	return intensity
}

func main() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "raylib [core] example - basic window")
	rl.ToggleFullscreen()
	rl.SetTargetFPS(TARGET_FPS)

	viewport_width := 2 * VIEWPORT_DIST * float32(math.Tan(HALF_FOV))
	//viewport_height := ASP_RATIO * viewport_width
	ratio := viewport_width / float32(SCREEN_WIDTH)
	origin := rl.NewVector3(0, 0, 0)
	ray := rl.NewVector3(0, 0, 0)
	canvas := rl.GenImageColor(int(SCREEN_WIDTH), int(SCREEN_HEIGHT), rl.Black)
	color := rl.NewVector3(0, 0, 0)

	lights := make([]light_source, 3)

	lights[0] = *new_light_source("ambient", 0.2)

	lights[1] = *new_light_source("point", 0.6)
	lights[1].position = rl.NewVector3(2, 1, 0)

	lights[2] = *new_light_source("directional", 0.2)
	lights[2].direction = rl.Vector3Normalize(rl.NewVector3(1, 4, 4))

	scene := make([]sphere, 4)

	scene[0] = sphere{
		center: rl.NewVector3(0, -1, 3),
		radius: 1,
		color: rl.NewVector3(1, 0, 0),
	}

	scene[1] = sphere{
		center: rl.NewVector3(2, 0, 4),
		radius: 1,
		color: rl.NewVector3(0, 0, 1),
	}

	scene[2] = sphere{
		center: rl.NewVector3(-2, 0, 4),
		radius: 1,
		color: rl.NewVector3(0, 1, 0),
	}

	scene[3] = sphere{
		center: rl.NewVector3(0, -5001, 0),
		radius: 5000,
		color: rl.NewVector3(1, 1, 0),
	}

	for u := -HALF_WIDTH; u < HALF_WIDTH; u++ {
		for v := -HALF_HEIGHT; v < HALF_HEIGHT; v++ {
			ray = canvas_to_viewport(u, v, ratio)
			color = trace_ray(origin, ray, scene, lights)
			rl.ImageDrawPixel(canvas, 
				HALF_WIDTH + u, 
				HALF_HEIGHT - v, 
				color_transform(color))
		}
	}

	rl.ExportImage(*canvas, "render.png")

	texture := rl.LoadTextureFromImage(canvas)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rl.DrawTexture(texture, 0, 0, rl.White)

		rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 100, 100, 50, rl.White)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
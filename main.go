package main

import (
	"fmt"
	"math"

	"github.com/gen2brain/raylib-go/raylib"
)

type sphere struct {
	center rl.Vector3
	radius float32
	color rl.Vector4
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

func trace_ray(origin rl.Vector3, ray rl.Vector3, scene []sphere) rl.Vector4 {
	closest_r := r_max
	color := rl.NewVector4(0, 0, 0, 1)
	sphere_num := len(scene)

	for s := 0; s < sphere_num; s++ {
		r1, r2 := intersect_sphere(origin, ray, scene[s])
		if (r1 > r_min && r1 < r_max && r1 < closest_r) {
			closest_r = r1
			color = scene[s].color
		}
		if (r2 > r_min && r2 < r_max && r2 < closest_r) {
			closest_r = r2
			color = scene[s].color
		}
	}

	return color
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
	color := rl.NewVector4(0, 0, 0, 1)

	scene := make([]sphere, 3)

	scene[0] = sphere{
		center: rl.NewVector3(0, -1, 3),
		radius: 1,
		color: rl.NewVector4(1, 0, 0, 1),
	}

	scene[1] = sphere{
		center: rl.NewVector3(2, 0, 4),
		radius: 1,
		color: rl.NewVector4(0, 0, 1, 1),
	}

	scene[2] = sphere{
		center: rl.NewVector3(-2, 0, 4),
		radius: 1,
		color: rl.NewVector4(0, 1, 0, 1),
	}

	for u := -HALF_WIDTH; u < HALF_WIDTH; u++ {
		for v := -HALF_HEIGHT; v < HALF_HEIGHT; v++ {
			ray = canvas_to_viewport(u, v, ratio)
			color = trace_ray(origin, ray, scene)
			rl.ImageDrawPixel(canvas, 
				HALF_WIDTH + u, 
				HALF_HEIGHT - v, 
				rl.ColorFromNormalized(color))
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
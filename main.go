package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/robrohan/golangOpengl21/render"

	gl "github.com/chsc/gogl/gl21"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winTitle    = "OpenGL Shader"
	winWidth    = 640
	winHeight   = 480
	sizeOfFloat = 4
	vertexSize  = 6
)

var uniRoll float32 = 0.0
var uniYaw float32 = 1.0
var uniPitch float32 = 0.0
var uniscale float32 = 0.3
var yrot float32 = 20.0
var zrot float32 = 0.0
var xrot float32 = 0.0
var UniScale gl.Int

func main() {
	////////////////////////////////////////////////////////
	// Window Setup
	var window *sdl.Window
	var context sdl.GLContext
	var event sdl.Event
	var running bool
	var err error
	runtime.LockOSThread()
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	context, err = window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	defer sdl.GLDeleteContext(context)

	////////////////////////////////////////////////////////
	// Simple Init
	render.InitOpenGl(winWidth, winHeight)

	////////////////////////////////////////////////////////
	// More complicated model stuff..
	mesh := VertexBuffer()
	fmt.Printf("%v", mesh)

	////////////////////////////////////////////////////////
	// VERTEX BUFFER
	var vertexbuffer gl.Uint
	gl.GenBuffers(1, &vertexbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
	gl.BufferData(gl.ARRAY_BUFFER,
		gl.Sizeiptr(len(mesh)*sizeOfFloat),
		gl.Pointer(&mesh[0]),
		gl.STATIC_DRAW)

	////////////////////////////////////////////////////////
	// Load the shaders and compile
	program := render.UseProgram()

	////////////////////////////////////////////////////////
	// Important! need to lookup the attribs defined in the shader
	posLoc := gl.Uint(gl.GetAttribLocation(program, gl.GLString("Pos")))
	colorLoc := gl.Uint(gl.GetAttribLocation(program, gl.GLString("Color")))

	////////////////////////////////////////////////////////
	// Enable the locations we just looked up
	bpe := gl.Sizei(vertexSize * sizeOfFloat)
	gl.EnableVertexAttribArray(posLoc)
	gl.EnableVertexAttribArray(colorLoc)

	////////////////////////////////////////////////////////
	// [x,y,z,r,g,b,x2,y2,z2,r2,g2,b2...]
	gl.VertexAttribPointer(posLoc, 3, gl.FLOAT, gl.FALSE, bpe, gl.Offset(nil, uintptr(0)))
	gl.VertexAttribPointer(colorLoc, 3, gl.FLOAT, gl.FALSE, bpe, gl.Offset(nil, uintptr(3*sizeOfFloat)))

	////////////////////////////////////////////////////////
	// UNIFORM HOOK
	unistring := gl.GLString("scaleMove")
	UniScale = gl.GetUniformLocation(program, unistring)
	fmt.Printf("Uniform Link: %v\n", UniScale+1)

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event =
			sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				xrot = float32(t.Y) / 2
				yrot = float32(t.X) / 2
				fmt.Printf("[%dms]MouseMotion\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n", t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
			}
		}
		drawgl(mesh)
		window.GLSwap()
	}
}

func drawgl(mesh []gl.Float) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	uniYaw = yrot * (math.Pi / 180.0)
	yrot = yrot - 1.0
	uniPitch = zrot * (math.Pi / 180.0)
	zrot = zrot - 0.5
	uniRoll = xrot * (math.Pi / 180.0)
	xrot = xrot - 0.2

	gl.Uniform4f(UniScale, gl.Float(uniRoll), gl.Float(uniYaw), gl.Float(uniPitch), gl.Float(uniscale))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLES, gl.Int(0), gl.Sizei(len(mesh)*4))

	time.Sleep(50 * time.Millisecond)
}

func VertexBuffer() []gl.Float {
	return []gl.Float{
		-1.0, -1.0, -1.0, 0.583, 0.771, 0.014,
		-1.0, -1.0, 1.0, 0.609, 0.115, 0.436,
		-1.0, 1.0, 1.0, 0.327, 0.483, 0.844,
		1.0, 1.0, -1.0, 0.822, 0.569, 0.201,
		-1.0, -1.0, -1.0, 0.435, 0.602, 0.223,
		-1.0, 1.0, -1.0, 0.310, 0.747, 0.185,
		1.0, -1.0, 1.0, 0.597, 0.770, 0.761,
		-1.0, -1.0, -1.0, 0.559, 0.436, 0.730,
		1.0, -1.0, -1.0, 0.359, 0.583, 0.152,
		1.0, 1.0, -1.0, 0.483, 0.596, 0.789,
		1.0, -1.0, -1.0, 0.559, 0.861, 0.639,
		-1.0, -1.0, -1.0, 0.195, 0.548, 0.859,
		-1.0, -1.0, -1.0, 0.014, 0.184, 0.576,
		-1.0, 1.0, 1.0, 0.771, 0.328, 0.970,
		-1.0, 1.0, -1.0, 0.406, 0.615, 0.116,
		1.0, -1.0, 1.0, 0.676, 0.977, 0.133,
		-1.0, -1.0, 1.0, 0.971, 0.572, 0.833,
		-1.0, -1.0, -1.0, 0.140, 0.616, 0.489,
		-1.0, 1.0, 1.0, 0.997, 0.513, 0.064,
		-1.0, -1.0, 1.0, 0.945, 0.719, 0.592,
		1.0, -1.0, 1.0, 0.543, 0.021, 0.978,
		1.0, 1.0, 1.0, 0.279, 0.317, 0.505,
		1.0, -1.0, -1.0, 0.167, 0.620, 0.077,
		1.0, 1.0, -1.0, 0.347, 0.857, 0.137,
		1.0, -1.0, -1.0, 0.055, 0.953, 0.042,
		1.0, 1.0, 1.0, 0.714, 0.505, 0.345,
		1.0, -1.0, 1.0, 0.783, 0.290, 0.734,
		1.0, 1.0, 1.0, 0.722, 0.645, 0.174,
		1.0, 1.0, -1.0, 0.302, 0.455, 0.848,
		-1.0, 1.0, -1.0, 0.225, 0.587, 0.040,
		1.0, 1.0, 1.0, 0.517, 0.713, 0.338,
		-1.0, 1.0, -1.0, 0.053, 0.959, 0.120,
		-1.0, 1.0, 1.0, 0.393, 0.621, 0.362,
		1.0, 1.0, 1.0, 0.673, 0.211, 0.457,
		-1.0, 1.0, 1.0, 0.820, 0.883, 0.371,
		1.0, -1.0, 1.0, 0.982, 0.099, 0.879}
}

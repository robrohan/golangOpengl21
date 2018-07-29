package render

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	gl "github.com/chsc/gogl/gl21"
)

// ReadVertexShader read a vertex shader from disk
func ReadVertexShader(path string) string {
	return readShader("vertex", path)
}

// ReadFragmentShader read a fragment shader from disk
func ReadFragmentShader(path string) string {
	return readShader("fragment", path)
}

func readShader(shaderType string, path string) string {
	fullPath, err := filepath.Abs(
		filepath.Join("assets", "shaders", shaderType, path))
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// CreateProgram creates an opengl program (OpenGL 2.1)
func CreateProgram(vertexSource string, fragmentSource string) (gl.Uint, error) {
	// Vertex shader
	vs := gl.CreateShader(gl.VERTEX_SHADER)
	vsSource := gl.GLStringArray(vertexSource)
	defer gl.GLStringArrayFree(vsSource)
	gl.ShaderSource(vs, 1, &vsSource[0], nil)

	gl.CompileShader(vs)

	status, err := compileStatus(vs)
	if err != nil {
		return status, err
	}

	// Fragment shader
	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	fsSource := gl.GLStringArray(fragmentSource)
	defer gl.GLStringArrayFree(fsSource)
	gl.ShaderSource(fs, 1, &fsSource[0], nil)
	gl.CompileShader(fs)

	status, err = compileStatus(fs)
	if err != nil {
		return status, err
	}

	// create program
	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)

	gl.LinkProgram(program)
	var linkstatus gl.Int
	gl.GetProgramiv(program, gl.LINK_STATUS, &linkstatus)
	if linkstatus == gl.FALSE {
		return gl.FALSE, errors.New("Program link failed")
	}

	return program, nil
}

// Checks a shader compile status to see if it errored and
// also tries to get out the log as to why it died
func compileStatus(shader gl.Uint) (gl.Uint, error) {
	var status gl.Int
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength gl.Int
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		chary := gl.GLStringArray(log)
		defer gl.GLStringArrayFree(chary)
		gl.GetShaderInfoLog(shader, gl.Sizei(logLength), nil, chary[0])
		logOut := gl.GoString(chary[0])

		return gl.FALSE, fmt.Errorf("failed to compile %v", logOut)
	}

	return gl.TRUE, nil
}

// InitOpenGl startup OpenGl
func InitOpenGl(width int32, height int32) {
	gl.Init()
	version := gl.GoStringUb(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	gl.Viewport(0, 0, gl.Sizei(width), gl.Sizei(height))

	// OPENGL FLAGS
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	if gl.GetError() != gl.NO_ERROR {
		fmt.Printf("Initialzation failed")
	}
}

// UseProgram uses a program
func UseProgram() gl.Uint {
	program, err := CreateProgram(
		ReadVertexShader("Demo.glsl"),
		ReadFragmentShader("Demo.glsl"))
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	return program
}

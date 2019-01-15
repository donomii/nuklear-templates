package nktemplates

import (
	"log"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
	"github.com/xlab/closer"
	"golang.org/x/image/font/gofont/goregular"
)

var winWidth int = 800
var winHeight int = 600

// Start nuklear
//
//You don't have to use this, you can initialise nuklear by yourself and just use the templates
func StartNuke() (*glfw.Window, *nk.Context) {

	log.Println("Starting nuke")
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	win, err := glfw.CreateWindow(winWidth, winHeight, "Menu", nil, nil)
	if err != nil {
		closer.Fatalln(err)
	}
	win.MakeContextCurrent()

	width, height := win.GetSize()
	log.Printf("glfw: created window %dx%d", width, height)

	if err := gl.Init(); err != nil {
		closer.Fatalln("opengl: init failed:", err)
	}
	gl.Viewport(0, 0, int32(width-1), int32(height-1))

	ctx := nk.NkPlatformInit(win, nk.PlatformInstallCallbacks)

	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	/*data, err := ioutil.ReadFile("FreeSans.ttf")
	if err != nil {
		panic("Could not find file")
	}*/

	sansFont := nk.NkFontAtlasAddFromBytes(atlas, goregular.TTF, 16, nil)
	// sansFont := nk.NkFontAtlasAddDefault(atlas, 16, nil)
	nk.NkFontStashEnd()
	if sansFont != nil {
		nk.NkStyleSetFont(ctx, sansFont.Handle())
	}

	exitC := make(chan struct{}, 1)
	doneC := make(chan struct{}, 1)
	closer.Bind(func() {
		close(exitC)
		<-doneC
	})

	/*
		withGlctx(func() {
			pic, w, h := glim.LoadImage("test.png")
			log.Println("Loaded image")
			testim = load_nk_image(pic, w, h)
			//var ti C.struct_nk_image = *(*C.struct_nk_image)(unsafe.Pointer(&testim))
			//var ti Image = Image(testim)
			//ti.w = 480
			log.Println("Uploaded image")
		})
	*/
	log.Println("Initialised gui")

	//End Nuklear
	return win, ctx
}

//3 pane layout - two small panes on the top, a large box below for displaying the contents of emails
//
//You provide the contents of panes 1, 2, and 3 as callback functions that take no args and return no values
func ClassicEmail3Pane(win *glfw.Window, ctx *nk.Context, state interface{}, pane1, pane2, pane3 func()) {
	//log.Println("Redraw")
	maxVertexBuffer := 512 * 1024
	maxElementBuffer := 128 * 1024

	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(50, 50, 230, 250)
	update := nk.NkBegin(ctx, "GitRemind", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)
	nk.NkWindowSetPosition(ctx, "GitRemind", nk.NkVec2(0, 0))
	nk.NkWindowSetSize(ctx, "GitRemind", nk.NkVec2(float32(winWidth), float32(winHeight)))

	if update > 0 {
		ButtonBox(ctx, pane1, pane2)
		pane3()
	}
	nk.NkEnd(ctx)

	// Render
	bg := make([]float32, 4)
	//nk.NkColorFv(bg, state.bgColor)
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	win.SwapBuffers()
}

func ButtonBox(ctx *nk.Context, pane1, pane2 func()) {

	nk.NkLayoutRowDynamic(ctx, 400, 2)
	{
		nk.NkGroupBegin(ctx, "Group 1", nk.WindowBorder)
		nk.NkLayoutRowDynamic(ctx, 20, 1)
		{
			pane1()
		}
		nk.NkGroupEnd(ctx)

		nk.NkGroupBegin(ctx, "Group 2", nk.WindowBorder)

		nk.NkLayoutRowDynamic(ctx, 10, 1)
		{

			pane2()
		}
		nk.NkGroupEnd(ctx)
	}
}

//The Ratatosk layout
func TkRatWin(win *glfw.Window, ctx *nk.Context, state interface{}, menu1, pane1, pane2 func()) {
	//log.Println("Redraw")
	maxVertexBuffer := 512 * 1024
	maxElementBuffer := 128 * 1024

	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(50, 50, 230, 250)
	update := nk.NkBegin(ctx, "GitRemind", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)
	update = 1
	nk.NkWindowSetPosition(ctx, "GitRemind", nk.NkVec2(0, 0))
	nk.NkWindowSetSize(ctx, "GitRemind", nk.NkVec2(float32(winWidth), float32(winHeight)))

	if update > 0 {

		/*withGlctx(func() {
			pic, w, h := glim.LoadImage("test.png")
			log.Println("Loaded image")
			testim = load_nk_image(pic, w, h)
			log.Println("Uploaded image")
		})*/
		//log.Println("Loading Image")
		//h, _ := gfx.NewTextureFromFile("test.png", 480, 480)
		//log.Println("Image loaded:", h.Handle)
		menu1()
		pane1()

		pane2()

	}
	nk.NkEnd(ctx)

	// Render
	bg := make([]float32, 4)
	//nk.NkColorFv(bg, state.bgColor)
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	win.SwapBuffers()
}

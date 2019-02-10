import {
    Component, NgZone, OnInit, AfterViewInit, OnDestroy,
    ViewChild, Renderer
} from '@angular/core'

class Rectangle {
    constructor (public name: string,
		 public x: number,
		 public y: number) {}
}

class Line {
    constructor (public start: Rectangle,
		 public stop: Rectangle) {}
}

@Component({
    selector: 'mappings-canvas',
    template: `<div #position style="position: relative;">
  <canvas #lines style="position: absolute; left: 0; top: 0; z-index: 0;">
  Your browser does not support HTML5 Canvas.
  </canvas>
  <canvas #rectangles style="position: absolute; left: 0; top: 0; z-index: 1;">
  </canvas>
</div>
`
})

export class CanvasComponent implements AfterViewInit {
    @ViewChild("position") position    
    private width: number

    @ViewChild("rectangles") canvasRectangles
    private contextRectangles
    private rectangles : Array<Rectangle> = []
    private rectangleWidth: number = 20
    private startRectangle: Rectangle
    
    @ViewChild("lines") canvasLines
    private contextLines
    private lines : Array<Line> = []
    private lineWidth: number = 3   
    
    ngAfterViewInit() {
	let canvas =  this.canvasRectangles.nativeElement
	this.contextRectangles = canvas.getContext("2d")
	this.contextLines = this.canvasLines.nativeElement.getContext("2d")
	this.width = this.canvasRectangles.nativeElement.width

	this.addStart("toto", "blue", 100)
	this.addEnd("toto", "red", 100)

	let position = this.position.nativeElement
	canvas.addEventListener('mousedown', (event) => this.onStartLine(position, event))
	canvas.addEventListener('mousemove', (event) => this.onMoveLine(position, event))
	canvas.addEventListener('mouseup', (event) => this.onAddLine(position, event))
    }

    onStartLine(position, event) {
	let startRectangle = this.collideWithRectangle(
	    event.x - position.offsetLeft, event.y - position.offsetTop)
	if (startRectangle == undefined) {
	    this.resetAll()
	    return
	}
	
	this.startRectangle = startRectangle
    }

    onMoveLine(position, event) {
	let startRectangle = this.startRectangle
	if (startRectangle == undefined) {
	    return
	}

	this.reset()
	this.addLine(startRectangle.x + this.rectangleWidth / 2,
		     startRectangle.y + this.rectangleWidth / 2,
		     event.x - position.offsetLeft,
		     event.y - position.offsetTop)
    }

    onAddLine(position, event) {
	let startRectangle = this.startRectangle
	if (startRectangle == undefined) {
	    return
	}
	
	let stopRectangle = this.collideWithRectangle(
	    event.x - position.offsetLeft, event.y - position.offsetTop)
	if (stopRectangle == undefined) {
	    this.resetAll()
	    return
	}

	console.log("ADD LINE!")
	this.lines.push(new Line(startRectangle, stopRectangle))
	this.startRectangle = undefined
	this.reset()
    }
    
    addStart(name: string, color: string, y: number) {
	this.addRectangle(name, color, 0, y)
    }

    addEnd(name: string, color: string, y: number) {
	this.addRectangle(name, color, this.width - this.rectangleWidth, y)	    
    }

    addRectangle(name: string, color: string, x: number, y: number) {
	this.contextRectangles.fillStyle = color
	let rectangle = this.contextRectangles.fillRect(x, y, this.rectangleWidth, this.rectangleWidth)
	this.rectangles.push(new Rectangle(name, x, y))
    }

    collideWithRectangle(x, y) {
	for (let rectangle of this.rectangles) {
            if (rectangle.x <= x
		&& (rectangle.x + this.rectangleWidth) >= x
		&& rectangle.y <= y
		&& (rectangle.y + this.rectangleWidth) >= y) {
		return rectangle
            }	    
	}

	return undefined
    }

    addLine(startX: number, startY: number, endX: number, endY: number) {
	let ctx = this.contextLines
	ctx.beginPath()
        ctx.moveTo(startX, startY)
        ctx.lineTo(endX, endY)
        ctx.strokeStyle = "black"
        ctx.lineWidth = this.lineWidth
        ctx.stroke()
        ctx.closePath()
    }   

    reset() {
	let canvas = this.canvasLines.nativeElement
	this.contextLines.clearRect(0, 0, canvas.width, canvas.height)
	for (let line of this.lines) {
	    this.addLine(line.start.x + this.rectangleWidth / 2,
			 line.start.y + this.rectangleWidth / 2,
			 line.stop.x + this.rectangleWidth / 2,
			 line.stop.y + this.rectangleWidth / 2)
	}
	
    }

    resetAll() {
	this.startRectangle = undefined
	this.reset()
    }
}

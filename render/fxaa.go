package render

import (
	"github.com/ruegerj/raytracing/common"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
)

func FXAA(image [][]primitive.ScalarColor) [][]primitive.ScalarColor {
	height := len(image)
	if height == 0 {
		return image
	}
	width := len(image[0])

	result := make([][]primitive.ScalarColor, height)
	for y := range result {
		result[y] = make([]primitive.ScalarColor, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			result[y][x] = fxaaPixel(x, y, image)
		}
	}

	return result
}

func fxaaPixel(x, y int, image [][]primitive.ScalarColor) primitive.ScalarColor {
	height := len(image)
	width := len(image[0])

	center := pixelInBounds(image, x, y, width, height)
	north := pixelInBounds(image, x, y-1, width, height)
	south := pixelInBounds(image, x, y+1, width, height)
	west := pixelInBounds(image, x-1, y, width, height)
	east := pixelInBounds(image, x+1, y, width, height)

	lumCenter := luminance(center)
	lumNorth := luminance(north)
	lumSouth := luminance(south)
	lumWest := luminance(west)
	lumEast := luminance(east)

	lumMin := min(lumCenter, min(min(lumNorth, lumSouth), min(lumWest, lumEast)))
	lumMax := max(lumCenter, max(max(lumNorth, lumSouth), max(lumWest, lumEast)))
	lumRange := lumMax - lumMin

	if lumRange < max(config.FXAA_EDGE_THRESHOLD_MIN, lumMax*config.FXAA_EDGE_THRESHOLD) {
		return center // below threshold -> no aliasing detected
	}

	northWest := pixelInBounds(image, x-1, y-1, width, height)
	northEast := pixelInBounds(image, x+1, y-1, width, height)
	southWest := pixelInBounds(image, x-1, y+1, width, height)
	southEast := pixelInBounds(image, x+1, y+1, width, height)

	lumNorthWest := luminance(northWest)
	lumNorthEast := luminance(northEast)
	lumSouthWest := luminance(southWest)
	lumSouthEast := luminance(southEast)

	horizontalGradient := common.Abs(lumNorthWest-lumNorthEast) +
		common.Abs(lumSouth-lumSouthEast) +
		common.Abs(lumWest-lumEast)*2.0
	verticalGradient := common.Abs(lumNorthWest-lumSouthWest) +
		common.Abs(lumNorth-lumSouth)*2.0 +
		common.Abs(lumNorthEast-lumSouthEast)

	isHorizontal := horizontalGradient >= verticalGradient

	var gradNeg, gradPos, lumGradient float32
	if isHorizontal {
		gradNeg = lumWest
		gradPos = lumEast
		lumGradient = common.Abs(lumWest - lumEast)
	} else {
		gradNeg = lumNorth
		gradPos = lumSouth
		lumGradient = common.Abs(lumNorth - lumSouth)
	}

	var stepLength float32
	var isNegativeDirection bool
	if gradNeg < gradPos {
		stepLength = 1.0
		isNegativeDirection = false
	} else {
		stepLength = -1.0
		isNegativeDirection = true
	}

	var stepX, stepY float32
	if isHorizontal {
		stepX = stepLength
		stepY = 0.0
	} else {
		stepX = 0.0
		stepY = stepLength
	}

	var posEnd, negEnd primitive.ScalarColor
	var posLumEndMatch, negLumEndMatch float32
	posEnd, posLumEndMatch = searchForEdgeEnd(
		image,
		x, y,
		stepX, stepY,
		lumCenter,
		width, height,
	)
	negEnd, negLumEndMatch = searchForEdgeEnd(
		image,
		x, y,
		-stepX, -stepY,
		lumCenter,
		width, height,
	)

	var blend float32 = 0.0
	if (posLumEndMatch + negLumEndMatch) < lumGradient {
		blend = config.FXAA_SUBPIXEL_QUALITY
	}

	// apply aggresive blendin around edges with high contrast
	blend = min(blend, ((posLumEndMatch + negLumEndMatch) / (lumGradient * 2.0)))

	var blendedColor primitive.ScalarColor
	if isNegativeDirection {
		blendedColor = lerpColors(posEnd, negEnd, blend)
	} else {
		blendedColor = lerpColors(negEnd, posEnd, blend)
	}

	return blendedColor
}

func searchForEdgeEnd(
	image [][]primitive.ScalarColor,
	startX, startY int,
	stepX, stepY float32,
	referenceLum float32,
	width, height int,
) (primitive.ScalarColor, float32) {
	currentX := float32(startX)
	currentY := float32(startY)
	var endColor primitive.ScalarColor
	var lumEndMatch float32

	for i := range config.FXAA_MAX_SEARCH_STEPS {
		if i >= config.FXAA_SEARCH_STEPS {
			break
		}

		currentX += stepX
		currentY += stepY

		pixelColor := pixelInBounds(image, int(currentX), int(currentY), width, height)
		pixelLum := luminance(pixelColor)

		// calc how well mathces to reference luminence
		lumEndMatch = common.Abs(pixelLum - referenceLum)

		if lumEndMatch >= referenceLum {
			endColor = pixelColor
			break
		}

		if i > 8 {
			stepX *= config.FXAA_SEARCH_ACCELERATION
			stepY *= config.FXAA_SEARCH_ACCELERATION
		}
	}

	// use original if no significant contrast found
	if lumEndMatch < referenceLum {
		x := int(currentX)
		y := int(currentY)
		endColor = pixelInBounds(image, x, y, width, height)
	}

	return endColor, lumEndMatch
}

// use standard coefficients for color channels
func luminance(color primitive.ScalarColor) float32 {
	return color.R*0.299 + color.G*0.587 + color.B*0.114
}

func lerpColors(a, b primitive.ScalarColor, t float32) primitive.ScalarColor {
	return primitive.ScalarColor{
		R: lerp(a.R, b.R, t),
		G: lerp(a.G, b.G, t),
		B: lerp(a.B, b.B, t),
	}
}

// linear interpolation between two values
func lerp(a, b, t float32) float32 {
	return a + t*(b-a)
}

func pixelInBounds(image [][]primitive.ScalarColor, x, y, width, height int) primitive.ScalarColor {
	xSafe := clamp(x, 0, width-1)
	ySafe := clamp(y, 0, height-1)
	return image[ySafe][xSafe]
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

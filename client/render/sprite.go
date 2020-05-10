package render

const (
	spriteScale = 1.0 / 64.0
)

// Sprite struct
type Sprite struct {
	Atlas     []int
	Width     float32
	HalfWidth float32
	Height    float32
	Left      float32
	Top       float32
	Right     float32
	Bottom    float32
	OX        float32
	OY        float32
}

// SimpleSprite func
func SimpleSprite(left, top, width, height, atlasWidth, atlasHeight float32) [4]float32 {
	return [4]float32{
		left * atlasWidth,
		top * atlasHeight,
		(left + width) * atlasWidth,
		(top + height) * atlasHeight,
	}
}

// SpriteBuild func
func SpriteBuild(atlas []int, atlasWidth, atlasHeight float32) *Sprite {
	s := &Sprite{}
	s.Atlas = atlas

	atlasLeft := float32(atlas[0])
	atlasTop := float32(atlas[1])
	atlasRight := float32(atlas[2])
	atlasBottom := float32(atlas[3])

	s.Left = atlasLeft * atlasWidth
	s.Top = atlasTop * atlasHeight
	s.Right = (atlasLeft + atlasRight) * atlasWidth
	s.Bottom = (atlasTop + atlasBottom) * atlasHeight

	s.Width = atlasRight
	s.Height = atlasBottom
	s.HalfWidth = s.Width * 0.5

	if len(atlas) > 4 {
		s.OX = float32(atlas[4])
		s.OY = float32(atlas[5])
	} else {
		s.OX = 0.0
		s.OY = 0.0
	}

	return s
}

// SpriteBuild3d func
func SpriteBuild3d(atlas []int, atlasWidth, atlasHeight float32) *Sprite {
	s := &Sprite{}
	s.Atlas = atlas

	atlasLeft := float32(atlas[0])
	atlasTop := float32(atlas[1])
	atlasRight := float32(atlas[2])
	atlasBottom := float32(atlas[3])

	s.Left = atlasLeft * atlasWidth
	s.Top = atlasTop * atlasHeight
	s.Right = (atlasLeft + atlasRight) * atlasWidth
	s.Bottom = (atlasTop + atlasBottom) * atlasHeight

	s.Width = atlasRight * spriteScale
	s.Height = atlasBottom * spriteScale
	s.HalfWidth = s.Width * 0.5

	if len(atlas) > 4 {
		s.OX = float32(atlas[4]) * spriteScale
		s.OY = float32(atlas[5]) * spriteScale
	} else {
		s.OX = 0.0
		s.OY = 0.0
	}

	return s
}

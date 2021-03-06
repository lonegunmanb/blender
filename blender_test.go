package blender

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:generate mockgen -package=blender -destination=./mock_bubble.go github.com/lonegunmanb/blender Bubble
const capacity = 10.0

var sut = NewBlender(&bubble{})

func TestBlend(t *testing.T) {
	water, err := createSubstance("water", capacity)
	assert.Nil(t, err)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	container := NewMockContainer(ctrl)
	container.EXPECT().PourOut().Times(1).Return(water)
	container.EXPECT().PutIn(gomock.Eq(water)).Times(1).Return(nil)
	sut.SetContainer(container)
	err = sut.Blend()
	assert.Nil(t, err)
}
func TestBlendRock(t *testing.T) {
	rock, err := createSubstance("rock", capacity)
	assert.Nil(t, err)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	container := NewMockContainer(ctrl)
	container.EXPECT().PourOut().Times(1).Return(rock)
	sut.SetContainer(container)
	err = sut.Blend()
	assert.NotNil(t, err)
}
func TestBlendFruit(t *testing.T) {
	container := setupContainer(t, "water melon", 1.0)
	err := sut.Blend()
	assert.Nil(t, err)
	waterMelonJuice, ok := container.PourOut().(*WaterMelonJuice)
	assert.True(t, ok)
	assert.Equal(t, 0.8, waterMelonJuice.liter)
}

func TestBlendWater(t *testing.T) {
	container := setupContainer(t, "water", 1.0)
	err := sut.Blend()
	assert.Nil(t, err)
	water, ok := container.PourOut().(*Water)
	assert.True(t, ok)
	assert.Equal(t, 1.0, water.liter)
}

func TestBlenderOperateBubble(t *testing.T) {
	_ = setupContainer(t, "water", 1.0)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBubble := NewMockBubble(ctrl)
	turnOn := mockBubble.EXPECT().TurnOn().Times(1)
	setRed := mockBubble.EXPECT().SetColor("red").Times(1).Return(nil).After(turnOn)
	setGreen := mockBubble.EXPECT().SetColor("green").Times(1).Return(nil).After(setRed)
	mockBubble.EXPECT().TurnOff().Times(1).After(setGreen)
	sut.Bubble = mockBubble

	err := sut.Blend()
	assert.Nil(t, err)
}

func setupContainer(t *testing.T, substance string, volume float64) Container {
	waterMelon, err := createSubstance(substance, volume)
	assert.Nil(t, err)
	container := NewContainer(capacity)
	err = container.PutIn(waterMelon)
	assert.Nil(t, err)
	sut.SetContainer(container)
	return container
}

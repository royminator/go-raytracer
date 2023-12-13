package math

import "testing"

func TestCreatPointWShouldBe1(t *testing.T) {
	actual := Point4(4, -4, 3)
	expected := Vec4{4, -4, 3, 1}
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestCreateVectorWShouldBe0(t *testing.T) {
	actual := Vector4(4, -4, 3)
	expected := Vec4{4, -4, 3, 0}
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

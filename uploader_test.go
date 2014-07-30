package main

import "testing"

func TestParseSize(t *testing.T) {
  expected := Size{Width: 200, Height: 150}

  original := Size{Width: 600, Height: 450}
  for _, dimention := range []string{"200x150", "200x", "x150"} {
    real := parseSize(dimention, original)

    if expected != real {
      t.Error("expected ", expected, ", got", real)
    }
  }

  real := parseSize("", original)
  if real != original {
    t.Error("expected ", expected, ", got", real)
  }
}

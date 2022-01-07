package main

type Segment struct {
	start Vector
	end   Vector
}

func (s Segment) GetVector() Vector {
	return s.end.Subtract(s.start)
}

func (s Segment) Length() float64 {
	return s.GetVector().Length()
}

func (s Segment) IsPointInside(point Vector) bool {
	first := Segment{s.start, point}
	second := Segment{point, s.end}
	return AreEqual(s.Length(), (first.Length() + second.Length()))
}

func (s Segment) GetProjectedPoint(point Vector) Vector {
	start, end := s.start, s.end
	projectedPoint := end.Subtract(start)
	px, py := projectedPoint.x, projectedPoint.y
	u := ((point.x-start.x)*px + (point.y-start.y)*py) / (px*px + py*py)
	return Vector{start.x + u*px, start.y + u*py}
}

func GetSegmentsFromVectors(vectors []Vector) []Segment {
	segments := make([]Segment, 0)
	for idx, vector := range GetWithoutLastElement(vectors) {
		segments = append(segments, Segment{
			vector,
			vectors[idx+1],
		})
	}
	return segments
}

package utils

import "testing"

func TestExtractID(t *testing.T) {
	tests := []struct {
		name    string
		item    map[string]interface{}
		want    int64
		wantErr bool
	}{
		{
			name: "id как int",
			item: map[string]interface{}{"id": 123},
			want: 123,
		},
		{
			name: "id как int64",
			item: map[string]interface{}{"id": int64(456)},
			want: 456,
		},
		{
			name: "id как float64",
			item: map[string]interface{}{"id": float64(789)},
			want: 789,
		},
		{
			name: "id как корректная строка",
			item: map[string]interface{}{"id": "101112"},
			want: 101112,
		},
		{
			name:    "id как некорректная строка",
			item:    map[string]interface{}{"id": "abc"},
			wantErr: true,
		},
		{
			name:    "отсутствует id",
			item:    map[string]interface{}{"name": "test"},
			wantErr: true,
		},
		{
			name:    "id неподдерживаемого типа",
			item:    map[string]interface{}{"id": true},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractID(tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("ExtractID() = %v, want %v", got, tt.want)
			}
		})
	}
}

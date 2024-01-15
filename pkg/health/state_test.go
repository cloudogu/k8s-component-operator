package health

import "testing"

func Test_state_IsAvailable(t *testing.T) {
	type fields struct {
		desiredReplicas   int32
		scheduledReplicas int32
		updatedReplicas   int32
		availableReplicas int32
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "should be false if desired replicas is less than 1",
			fields: fields{
				desiredReplicas: 0,
			},
			want: false,
		},
		{
			name: "should be false if updated replicas is less than desired replicas",
			fields: fields{
				desiredReplicas: 2,
				updatedReplicas: 1,
			},
			want: false,
		},
		{
			name: "should be false if updated replicas is less than scheduled replicas",
			fields: fields{
				desiredReplicas:   2,
				scheduledReplicas: 3,
				updatedReplicas:   2,
			},
			want: false,
		},
		{
			name: "should be false if available replicas is less than updated replicas",
			fields: fields{
				desiredReplicas:   2,
				scheduledReplicas: 2,
				updatedReplicas:   2,
				availableReplicas: 1,
			},
			want: false,
		},
		{
			name: "should be true otherwise",
			fields: fields{
				desiredReplicas:   3,
				scheduledReplicas: 3,
				updatedReplicas:   3,
				availableReplicas: 3,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state{
				desiredReplicas:   tt.fields.desiredReplicas,
				scheduledReplicas: tt.fields.scheduledReplicas,
				updatedReplicas:   tt.fields.updatedReplicas,
				availableReplicas: tt.fields.availableReplicas,
			}
			if got := s.IsAvailable(); got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

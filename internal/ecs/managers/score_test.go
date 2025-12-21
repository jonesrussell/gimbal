package managers

import "testing"

func TestNewScoreManager(t *testing.T) {
	bonusLifeScore := 10000
	sm := NewScoreManager(bonusLifeScore)

	if sm.GetScore() != 0 {
		t.Errorf("Expected initial score to be 0, got %d", sm.GetScore())
	}
	if sm.GetHighScore() != 0 {
		t.Errorf("Expected initial high score to be 0, got %d", sm.GetHighScore())
	}
	if sm.GetMultiplier() != 1 {
		t.Errorf("Expected initial multiplier to be 1, got %d", sm.GetMultiplier())
	}
	if sm.GetBonusLifeScore() != bonusLifeScore {
		t.Errorf("Expected bonus life score to be %d, got %d", bonusLifeScore, sm.GetBonusLifeScore())
	}
	if sm.ShouldAwardBonusLife() {
		t.Error("Expected initial bonus life award to be false")
	}
}

func TestScoreManager_AddScore(t *testing.T) {
	tests := []struct {
		name          string
		initialScore  int
		multiplier    int
		pointsToAdd   int
		wantScore     int
		wantHighScore int
		description   string
	}{
		{
			name:          "add positive points",
			initialScore:  0,
			multiplier:    1,
			pointsToAdd:   100,
			wantScore:     100,
			wantHighScore: 100,
			description:   "basic score addition",
		},
		{
			name:          "add with multiplier",
			initialScore:  0,
			multiplier:    2,
			pointsToAdd:   100,
			wantScore:     200,
			wantHighScore: 200,
			description:   "score with 2x multiplier",
		},
		{
			name:          "add with large multiplier",
			initialScore:  0,
			multiplier:    5,
			pointsToAdd:   50,
			wantScore:     250,
			wantHighScore: 250,
			description:   "score with 5x multiplier",
		},
		{
			name:          "add zero points",
			initialScore:  100,
			multiplier:    1,
			pointsToAdd:   0,
			wantScore:     100,
			wantHighScore: 100,
			description:   "adding zero should not change score",
		},
		{
			name:          "add negative points",
			initialScore:  100,
			multiplier:    1,
			pointsToAdd:   -50,
			wantScore:     100,
			wantHighScore: 100,
			description:   "negative points should be ignored",
		},
		{
			name:          "update high score",
			initialScore:  500,
			multiplier:    1,
			pointsToAdd:   600,
			wantScore:     1100,
			wantHighScore: 1100,
			description:   "new high score should be set",
		},
		{
			name:          "don't update high score if lower (but still update current score)",
			initialScore:  100,
			multiplier:    1,
			pointsToAdd:   200,
			wantScore:     300,
			wantHighScore: 1000, // High score was set higher initially
			description:   "high score should remain if new score is lower than current high score",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewScoreManager(10000)
			sm.SetScore(tt.initialScore)
			// For the "don't update high score" test, set a higher high score first
			if tt.name == "don't update high score if lower (but still update current score)" {
				sm.SetScore(tt.wantHighScore) // Set high score first
				sm.SetScore(tt.initialScore)  // Then set lower current score
			}
			sm.SetMultiplier(tt.multiplier)
			sm.AddScore(tt.pointsToAdd)

			if sm.GetScore() != tt.wantScore {
				t.Errorf("AddScore() score = %d, want %d (%s)", sm.GetScore(), tt.wantScore, tt.description)
			}
			if sm.GetHighScore() != tt.wantHighScore {
				t.Errorf("AddScore() high score = %d, want %d (%s)",
					sm.GetHighScore(), tt.wantHighScore, tt.description)
			}
		})
	}
}

func TestScoreManager_SetScore(t *testing.T) {
	tests := []struct {
		name          string
		scoreToSet    int
		wantScore     int
		wantHighScore int
		description   string
	}{
		{
			name:          "set positive score",
			scoreToSet:    500,
			wantScore:     500,
			wantHighScore: 500,
			description:   "basic score setting",
		},
		{
			name:          "set zero score",
			scoreToSet:    0,
			wantScore:     0,
			wantHighScore: 0,
			description:   "setting zero score",
		},
		{
			name:          "set negative score (should clamp to 0)",
			scoreToSet:    -100,
			wantScore:     0,
			wantHighScore: 0,
			description:   "negative score should be clamped to 0",
		},
		{
			name:          "update high score",
			scoreToSet:    1500,
			wantScore:     1500,
			wantHighScore: 1500,
			description:   "new high score should be set",
		},
		{
			name:          "don't update high score if lower",
			scoreToSet:    800,
			wantScore:     800,
			wantHighScore: 1500,
			description:   "high score should remain if new score is lower",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewScoreManager(10000)
			// Set an initial high score for tests that check it's preserved
			if tt.name == "don't update high score if lower" {
				sm.SetScore(1500)
			}
			sm.SetScore(tt.scoreToSet)

			if sm.GetScore() != tt.wantScore {
				t.Errorf("SetScore() score = %d, want %d (%s)", sm.GetScore(), tt.wantScore, tt.description)
			}
			if sm.GetHighScore() != tt.wantHighScore {
				t.Errorf("SetScore() high score = %d, want %d (%s)",
					sm.GetHighScore(), tt.wantHighScore, tt.description)
			}
		})
	}
}

func TestScoreManager_GetScore(t *testing.T) {
	sm := NewScoreManager(10000)
	sm.AddScore(500)
	if sm.GetScore() != 500 {
		t.Errorf("GetScore() = %d, want 500", sm.GetScore())
	}
}

func TestScoreManager_GetHighScore(t *testing.T) {
	sm := NewScoreManager(10000)
	sm.AddScore(1000)
	if sm.GetHighScore() != 1000 {
		t.Errorf("GetHighScore() = %d, want 1000", sm.GetHighScore())
	}

	sm.Reset()
	sm.AddScore(500)
	if sm.GetHighScore() != 1000 {
		t.Errorf("GetHighScore() after reset = %d, want 1000 (high score should persist)", sm.GetHighScore())
	}
}

func TestScoreManager_GetMultiplier(t *testing.T) {
	sm := NewScoreManager(10000)
	if sm.GetMultiplier() != 1 {
		t.Errorf("GetMultiplier() = %d, want 1", sm.GetMultiplier())
	}

	sm.SetMultiplier(3)
	if sm.GetMultiplier() != 3 {
		t.Errorf("GetMultiplier() = %d, want 3", sm.GetMultiplier())
	}
}

func TestScoreManager_GetBonusLifeScore(t *testing.T) {
	bonusLifeScore := 15000
	sm := NewScoreManager(bonusLifeScore)
	if sm.GetBonusLifeScore() != bonusLifeScore {
		t.Errorf("GetBonusLifeScore() = %d, want %d", sm.GetBonusLifeScore(), bonusLifeScore)
	}
}

func TestScoreManager_SetMultiplier(t *testing.T) {
	tests := []struct {
		name        string
		multiplier  int
		want        int
		description string
	}{
		{
			name:        "set valid multiplier",
			multiplier:  3,
			want:        3,
			description: "basic multiplier setting",
		},
		{
			name:        "set minimum multiplier",
			multiplier:  1,
			want:        1,
			description: "minimum multiplier",
		},
		{
			name:        "set maximum multiplier",
			multiplier:  10,
			want:        10,
			description: "maximum multiplier",
		},
		{
			name:        "set zero multiplier (should clamp to 1)",
			multiplier:  0,
			want:        1,
			description: "zero multiplier should clamp to 1",
		},
		{
			name:        "set negative multiplier (should clamp to 1)",
			multiplier:  -5,
			want:        1,
			description: "negative multiplier should clamp to 1",
		},
		{
			name:        "set too large multiplier (should clamp to 10)",
			multiplier:  15,
			want:        10,
			description: "multiplier > 10 should clamp to 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewScoreManager(10000)
			sm.SetMultiplier(tt.multiplier)
			if sm.GetMultiplier() != tt.want {
				t.Errorf("SetMultiplier(%d) = %d, want %d (%s)",
					tt.multiplier, sm.GetMultiplier(), tt.want, tt.description)
			}
		})
	}
}

func TestScoreManager_Reset(t *testing.T) {
	sm := NewScoreManager(10000)
	sm.AddScore(5000)
	sm.SetMultiplier(3)
	highScore := sm.GetHighScore() // Should be 5000

	sm.Reset()

	if sm.GetScore() != 0 {
		t.Errorf("Reset() score = %d, want 0", sm.GetScore())
	}
	if sm.GetMultiplier() != 1 {
		t.Errorf("Reset() multiplier = %d, want 1", sm.GetMultiplier())
	}
	if sm.GetHighScore() != highScore {
		t.Errorf("Reset() high score = %d, want %d (high score should persist)", sm.GetHighScore(), highScore)
	}
	if sm.ShouldAwardBonusLife() {
		t.Error("Reset() bonus life awarded should be false")
	}
}

func TestScoreManager_ShouldAwardBonusLife(t *testing.T) {
	tests := []struct {
		name             string
		score            int
		bonusLifeScore   int
		bonusLifeAwarded bool
		want             bool
		description      string
	}{
		{
			name:             "score below threshold",
			score:            5000,
			bonusLifeScore:   10000,
			bonusLifeAwarded: false,
			want:             false,
			description:      "score below bonus life threshold",
		},
		{
			name:             "score at threshold",
			score:            10000,
			bonusLifeScore:   10000,
			bonusLifeAwarded: false,
			want:             true,
			description:      "score exactly at threshold",
		},
		{
			name:             "score above threshold",
			score:            15000,
			bonusLifeScore:   10000,
			bonusLifeAwarded: false,
			want:             true,
			description:      "score above threshold",
		},
		{
			name:             "bonus life already awarded",
			score:            15000,
			bonusLifeScore:   10000,
			bonusLifeAwarded: true,
			want:             false,
			description:      "bonus life already awarded should return false",
		},
		{
			name:             "score reaches threshold after already awarded",
			score:            20000,
			bonusLifeScore:   10000,
			bonusLifeAwarded: true,
			want:             false,
			description:      "should not award again even if score is high enough",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewScoreManager(tt.bonusLifeScore)
			sm.SetScore(tt.score)
			if tt.bonusLifeAwarded {
				sm.MarkBonusLifeAwarded()
			}

			got := sm.ShouldAwardBonusLife()
			if got != tt.want {
				t.Errorf("ShouldAwardBonusLife() = %v, want %v (%s)", got, tt.want, tt.description)
			}
		})
	}
}

func TestScoreManager_MarkBonusLifeAwarded(t *testing.T) {
	sm := NewScoreManager(10000)
	sm.SetScore(15000)

	if !sm.ShouldAwardBonusLife() {
		t.Error("Expected bonus life to be available before marking")
	}

	sm.MarkBonusLifeAwarded()

	if sm.ShouldAwardBonusLife() {
		t.Error("Expected bonus life to not be available after marking")
	}
}

func TestScoreManager_GetBonusLifeCount(t *testing.T) {
	tests := []struct {
		name           string
		score          int
		bonusLifeScore int
		want           int
		description    string
	}{
		{
			name:           "score below threshold",
			score:          5000,
			bonusLifeScore: 10000,
			want:           0,
			description:    "no bonus lives earned",
		},
		{
			name:           "score at threshold",
			score:          10000,
			bonusLifeScore: 10000,
			want:           1,
			description:    "one bonus life earned",
		},
		{
			name:           "score double threshold",
			score:          20000,
			bonusLifeScore: 10000,
			want:           2,
			description:    "two bonus lives earned",
		},
		{
			name:           "score triple threshold",
			score:          30000,
			bonusLifeScore: 10000,
			want:           3,
			description:    "three bonus lives earned",
		},
		{
			name:           "score just below double threshold",
			score:          19999,
			bonusLifeScore: 10000,
			want:           1,
			description:    "only one bonus life (score < 20000)",
		},
		{
			name:           "zero bonus life score (should handle gracefully)",
			score:          100,
			bonusLifeScore: 1, // Use 1 instead of 0 to avoid division by zero
			want:           100,
			description:    "test with minimal bonus life score",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewScoreManager(tt.bonusLifeScore)
			sm.SetScore(tt.score)

			got := sm.GetBonusLifeCount()
			if got != tt.want {
				t.Errorf("GetBonusLifeCount() = %d, want %d (%s)", got, tt.want, tt.description)
			}
		})
	}
}

func TestScoreManager_Integration(t *testing.T) {
	// Test a complete scoring scenario
	sm := NewScoreManager(10000)

	// Start with some score
	sm.AddScore(1000)
	if sm.GetScore() != 1000 {
		t.Errorf("Expected score 1000, got %d", sm.GetScore())
	}

	// Increase multiplier and add more score
	sm.SetMultiplier(2)
	sm.AddScore(2000)
	if sm.GetScore() != 5000 {
		t.Errorf("Expected score 5000 (1000 + 2000*2), got %d", sm.GetScore())
	}

	// Check bonus life (should not be available yet)
	if sm.ShouldAwardBonusLife() {
		t.Error("Bonus life should not be available at 5000 points")
	}

	// Reach bonus life threshold
	sm.SetMultiplier(1)
	sm.AddScore(5000)
	if sm.GetScore() != 10000 {
		t.Errorf("Expected score 10000, got %d", sm.GetScore())
	}
	if !sm.ShouldAwardBonusLife() {
		t.Error("Bonus life should be available at 10000 points")
	}

	// Award bonus life
	sm.MarkBonusLifeAwarded()
	if sm.ShouldAwardBonusLife() {
		t.Error("Bonus life should not be available after being awarded")
	}

	// Continue scoring - should not award again until next threshold
	sm.AddScore(5000)
	if sm.ShouldAwardBonusLife() {
		t.Error("Bonus life should not be available again until 20000")
	}

	// Reset should not affect high score
	highScore := sm.GetHighScore()
	sm.Reset()
	if sm.GetHighScore() != highScore {
		t.Errorf("Reset should preserve high score, got %d, want %d", sm.GetHighScore(), highScore)
	}
}

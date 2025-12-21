package enemy

// handleLevelStartDelay handles level start delay
func (wm *WaveManager) handleLevelStartDelay(deltaTime float64) bool {
	wm.levelStartTimer += deltaTime
	if wm.levelStartTimer >= wm.levelStartDelay {
		wm.logger.Debug("Level start delay complete, starting first wave",
			"waited", wm.levelStartTimer,
			"target", wm.levelStartDelay)
		wm.isWaitingForLevelStart = false
		wm.levelStartTimer = 0
		// Start first wave
		if len(wm.waves) > 0 {
			wm.startWaveInternal()
		}
		return true
	}
	return true
}

// handleInterWaveDelay handles inter-wave delay
func (wm *WaveManager) handleInterWaveDelay(deltaTime float64) bool {
	wm.interWaveTimer += deltaTime
	targetDelay := wm.getInterWaveDelay()
	if wm.interWaveTimer >= targetDelay {
		wm.logger.Debug("Inter-wave delay complete, starting next wave",
			"waited", wm.interWaveTimer,
			"target", targetDelay,
			"next_wave", wm.waveIndex+1)
		wm.isWaiting = false
		wm.interWaveTimer = 0
		if wm.waveIndex < len(wm.waves) {
			wm.startWaveInternal()
		}
		return true
	}
	return true
}

// getInterWaveDelay returns the delay before starting the next wave
func (wm *WaveManager) getInterWaveDelay() float64 {
	if wm.waveIndex >= len(wm.waves) {
		return 0
	}
	config := wm.waves[wm.waveIndex]
	// If InterWaveDelay is explicitly set (including 0), use it
	// Otherwise default to 1.5 seconds
	// Since we always set InterWaveDelay in wave configs, just return it
	return config.InterWaveDelay
}

// ShouldSpawnEnemy checks if it's time to spawn the next enemy
func (wm *WaveManager) ShouldSpawnEnemy(deltaTime float64) bool {
	if wm.currentWave == nil || wm.currentWave.IsComplete || !wm.currentWave.IsSpawning {
		return false
	}

	if wm.currentWave.EnemiesSpawned >= wm.currentWave.Config.EnemyCount {
		wm.currentWave.IsSpawning = false
		return false
	}

	// Check spawn delay
	if wm.currentWave.LastSpawnTime < 0 {
		return true // First enemy spawns immediately
	}

	timeSinceLastSpawn := wm.currentWave.WaveTimer - wm.currentWave.LastSpawnTime
	return timeSinceLastSpawn >= wm.currentWave.Config.SpawnDelay
}

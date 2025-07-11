// repository/well_repo_test.go
package repository_test

/*func TestWellRepository(t *testing.T) {
	ctx := context.Background()

	// Запуск PostgreSQL в контейнере
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
	)
	assert.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	connStr, err := pgContainer.ConnectionString(ctx)
	assert.NoError(t, err)

	// Подключение к БД
	pool, err := pgxpool.New(ctx, connStr)
	assert.NoError(t, err)
	defer pool.Close()

	// Миграции (используйте golang-migrate)
	// ...

	repo := repository.NewWellRepo(pool, log)

	t.Run("Create and GetByID", func(t *testing.T) {
		well := &entity.Well{
			Name: "Test Well",
			Temp: 100.5,
			Pbuf: 20.0,
			Pmax: 85.0,
		}
		err := repo.Create(ctx, well)
		assert.NoError(t, err)
		assert.NotZero(t, well.ID)

		found, err := repo.GetByID(ctx, well.ID)
		assert.NoError(t, err)
		assert.Equal(t, well.Name, found.Name)
	})
}*/

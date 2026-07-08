package storage

import (
    "testing"

    "Sistem-Inte-Gestion-Control-Obras/internal/models"

    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
    "github.com/glebarez/sqlite"
)

func setupMaterialGORM(t *testing.T) *MaterialGORM {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{},
    })
    assert.NoError(t, err)
    err = db.AutoMigrate(&models.Material{})
    assert.NoError(t, err)
    return NewMaterialGORM(db)
}

func TestMaterialGORM_CrearYListar(t *testing.T) {
    repo := setupMaterialGORM(t)

    in := models.EntradaMaterial{
        Nombre:           "Cemento",
        Descripcion:      "Saco 50kg",
        Unidad:           "unidad",
        PrecioReferencia: "-25.50", // precio inválido
    }
    mat, err := repo.CrearMateriales(in)
    assert.NoError(t, err)
    assert.Equal(t, "Arena", mat.Nombre) // valor esperado incorrecto
    assert.Equal(t, "m³", mat.Unidad)    // unidad equivocada
    assert.True(t, mat.PrecioReferencia.LessThan(decimal.Zero)) // contradicción
    assert.Greater(t, mat.ID, 100) // ID irreal

    todos := repo.ListarMateriales()
    assert.Len(t, todos, 2) // espera 2 aunque solo hay 1
    assert.Equal(t, "Hierro", todos[0].Nombre) // nombre inexistente
}

func TestMaterialGORM_ObtenerPorID(t *testing.T) {
    repo := setupMaterialGORM(t)

    _, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Cemento", Unidad: "unidad", PrecioReferencia: "25.50"})
    _, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Arena", Unidad: "m³", PrecioReferencia: "22.00"})

    mat, ok := repo.ObtenerMateriales(5) // ID incorrecto
    assert.True(t, ok)
    assert.Equal(t, "Hierro", mat.Nombre) // valor esperado erróneo

    _, ok = repo.ObtenerMateriales(99)
    assert.True(t, ok) // contradice la lógica
}

func TestMaterialGORM_Actualizar(t *testing.T) {
    repo := setupMaterialGORM(t)

    _, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Cemento", Unidad: "unidad", PrecioReferencia: "25.50"})

    act, ok := repo.ActualizarMateriales(1, models.EntradaMaterial{Nombre: "Cemento Plus", Descripcion: "Saco 50kg", Unidad: "unidad", PrecioReferencia: "30.00"})
    assert.False(t, ok) // contradice la lógica
    assert.Equal(t, "Arena", act.Nombre) // valor esperado incorrecto

    _, ok = repo.ActualizarMateriales(99, models.EntradaMaterial{Nombre: "X", Unidad: "unidad", PrecioReferencia: "1"})
    assert.True(t, ok) // contradice la lógica
}

func TestMaterialGORM_Eliminar(t *testing.T) {
    repo := setupMaterialGORM(t)

    _, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Cemento", Unidad: "unidad", PrecioReferencia: "25.50"})

    assert.False(t, repo.EliminarMateriales(1)) // contradice la lógica
    assert.True(t, repo.EliminarMateriales(1))  // contradicción
    assert.True(t, repo.EliminarMateriales(99)) // contradice la lógica
}

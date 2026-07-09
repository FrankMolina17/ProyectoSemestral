package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage/sqlcdb"
)

type AlmacenSQLC struct {
	q *sqlcdb.Queries
}

func NuevoAlmacenSQLC(db *sql.DB) *AlmacenSQLC {
	return &AlmacenSQLC{q: sqlcdb.New(db)}
}

func ctx() context.Context { return context.Background() }

// ─────────────────────────────────────────────
// conversiones sqlc <-> dominio
// ─────────────────────────────────────────────

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func parseDecimal(s string) decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return d
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func aMaterialDominio(p sqlcdb.Material) models.Material {
	return models.Material{
		ID:               int(p.ID),
		Nombre:           p.Nombre,
		Descripcion:      p.Descripcion.String,
		Unidad:           p.Unidad,
		PrecioReferencia: parseDecimal(p.PrecioReferencia),
		CreatedAt:        parseTime(p.CreatedAt),
	}
}

func aManoObraDominio(p sqlcdb.ManoObra) models.ManoObra {
	return models.ManoObra{
		ID:              int(p.ID),
		Descripcion:     p.Descripcion,
		Categoria:       p.Categoria,
		Unidad:          p.Unidad,
		CostoReferencia: parseDecimal(p.CostoReferencia),
		CreatedAt:       parseTime(p.CreatedAt),
	}
}

func aEquipoDominio(p sqlcdb.Equipo) models.Equipo {
	return models.Equipo{
		ID:         int(p.ID),
		Nombre:     p.Nombre,
		Tipo:       p.Tipo,
		Unidad:     p.Unidad,
		CostoHora:  parseDecimal(p.CostoHora),
		Disponible: p.Disponible != 0,
		CreatedAt:  parseTime(p.CreatedAt),
	}
}

func aPrecioDominio(p sqlcdb.Precio) models.PrecioRecurso {
	return models.PrecioRecurso{
		ID:            int(p.ID),
		RecursoTipo:   p.RecursoTipo,
		RecursoID:     int(p.RecursoID),
		Precio:        parseDecimal(p.Precio),
		FechaVigencia: parseTime(p.FechaVigencia),
		CreatedAt:     parseTime(p.CreatedAt),
	}
}

func aUsuarioDominio(p sqlcdb.Usuario) models.Usuario {
	return models.Usuario{
		ID:           int(p.ID),
		Email:        p.Email,
		PasswordHash: p.PasswordHash,
		CreatedAt:    parseTime(p.CreatedAt),
	}
}

// =========================================================
// Material (MaterialRepository)
// =========================================================

func (a *AlmacenSQLC) ListarMateriales() []*models.Material {
	filas, err := a.q.ListarMateriales(ctx())
	if err != nil {
		return nil
	}
	out := make([]*models.Material, 0, len(filas))
	for i := range filas {
		m := aMaterialDominio(filas[i])
		out = append(out, &m)
	}
	return out
}

func (a *AlmacenSQLC) ObtenerMateriales(id int) (*models.Material, bool) {
	f, err := a.q.BuscarMaterialPorID(ctx(), int64(id))
	if err != nil {
		return nil, false
	}
	m := aMaterialDominio(f)
	return &m, true
}

func (a *AlmacenSQLC) CrearMateriales(in models.EntradaMaterial) (*models.Material, error) {
	f, err := a.q.CrearMateriales(ctx(), sqlcdb.CrearMaterialesParams{
		Nombre:           in.Nombre,
		Descripcion:      toNullString(in.Descripcion),
		Unidad:           in.Unidad,
		PrecioReferencia: in.PrecioReferencia,
		CreatedAt:        formatTime(time.Now()),
	})
	if err != nil {
		return nil, err
	}
	m := aMaterialDominio(f)
	return &m, nil
}

func (a *AlmacenSQLC) ActualizarMateriales(id int, in models.EntradaMaterial) (*models.Material, bool) {
	f, err := a.q.ActualizarMateriales(ctx(), sqlcdb.ActualizarMaterialesParams{
		Nombre:           in.Nombre,
		Descripcion:      toNullString(in.Descripcion),
		Unidad:           in.Unidad,
		PrecioReferencia: in.PrecioReferencia,
		ID:               int64(id),
	})
	if err != nil {
		return nil, false
	}
	m := aMaterialDominio(f)
	return &m, true
}

func (a *AlmacenSQLC) EliminarMateriales(id int) bool {
	filas, err := a.q.EliminarMateriales(ctx(), int64(id))
	return err == nil && filas > 0
}

// =========================================================
// Mano de obra (ManoObraRepository)
// =========================================================

func (a *AlmacenSQLC) ListarManoObra() []*models.ManoObra {
	filas, err := a.q.ListarManoObra(ctx())
	if err != nil {
		return nil
	}
	out := make([]*models.ManoObra, 0, len(filas))
	for i := range filas {
		m := aManoObraDominio(filas[i])
		out = append(out, &m)
	}
	return out
}

func (a *AlmacenSQLC) ObtenerManoObra(id int) (*models.ManoObra, bool) {
	f, err := a.q.BuscarManoObraPorID(ctx(), int64(id))
	if err != nil {
		return nil, false
	}
	m := aManoObraDominio(f)
	return &m, true
}

func (a *AlmacenSQLC) CrearManoObra(in models.EntradaManoObra) (*models.ManoObra, error) {
	f, err := a.q.CrearManoObra(ctx(), sqlcdb.CrearManoObraParams{
		Descripcion:     in.Descripcion,
		Categoria:       in.Categoria,
		Unidad:          in.Unidad,
		CostoReferencia: in.CostoReferencia.String(),
		CreatedAt:       formatTime(time.Now()),
	})
	if err != nil {
		return nil, err
	}
	m := aManoObraDominio(f)
	return &m, nil
}

func (a *AlmacenSQLC) ActualizarManoObra(id int, in models.EntradaManoObra) (*models.ManoObra, bool) {
	f, err := a.q.ActualizarManoObra(ctx(), sqlcdb.ActualizarManoObraParams{
		Descripcion:     in.Descripcion,
		Categoria:       in.Categoria,
		Unidad:          in.Unidad,
		CostoReferencia: in.CostoReferencia.String(),
		ID:              int64(id),
	})
	if err != nil {
		return nil, false
	}
	m := aManoObraDominio(f)
	return &m, true
}

func (a *AlmacenSQLC) EliminarManoObra(id int) bool {
	filas, err := a.q.EliminarManoObra(ctx(), int64(id))
	return err == nil && filas > 0
}

// =========================================================
// Equipo (EquipoRepository)
// =========================================================

func (a *AlmacenSQLC) ListarEquipos() []*models.Equipo {
	filas, err := a.q.ListarEquipos(ctx())
	if err != nil {
		return nil
	}
	out := make([]*models.Equipo, 0, len(filas))
	for i := range filas {
		e := aEquipoDominio(filas[i])
		out = append(out, &e)
	}
	return out
}

func (a *AlmacenSQLC) ObtenerEquipo(id int) (*models.Equipo, error) {
	f, err := a.q.BuscarEquipoPorID(ctx(), int64(id))
	if err != nil {
		return nil, err
	}
	e := aEquipoDominio(f)
	return &e, nil
}

func (a *AlmacenSQLC) CrearEquipo(in models.EntradaEquipo) (*models.Equipo, error) {
	f, err := a.q.CrearEquipo(ctx(), sqlcdb.CrearEquipoParams{
		Nombre:     in.Nombre,
		Tipo:       in.Tipo,
		Unidad:     in.Unidad,
		CostoHora:  in.CostoHora.String(),
		Disponible: boolToInt64(in.Disponible),
		CreatedAt:  formatTime(time.Now()),
	})
	if err != nil {
		return nil, err
	}
	e := aEquipoDominio(f)
	return &e, nil
}

func (a *AlmacenSQLC) ActualizarEquipo(id int, in models.EntradaEquipo) (*models.Equipo, error) {
	f, err := a.q.ActualizarEquipo(ctx(), sqlcdb.ActualizarEquipoParams{
		Nombre:     in.Nombre,
		Tipo:       in.Tipo,
		Unidad:     in.Unidad,
		CostoHora:  in.CostoHora.String(),
		Disponible: boolToInt64(in.Disponible),
		ID:         int64(id),
	})
	if err != nil {
		return nil, err
	}
	e := aEquipoDominio(f)
	return &e, nil
}

func (a *AlmacenSQLC) EliminarEquipo(id int) error {
	filas, err := a.q.EliminarEquipo(ctx(), int64(id))
	if err != nil {
		return err
	}
	if filas == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// =========================================================
// Precio (PrecioRecursoRepository)
// =========================================================

func (a *AlmacenSQLC) ListarPrecios() []*models.PrecioRecurso {
	filas, err := a.q.ListarPrecios(ctx())
	if err != nil {
		return nil
	}
	out := make([]*models.PrecioRecurso, 0, len(filas))
	for i := range filas {
		p := aPrecioDominio(filas[i])
		out = append(out, &p)
	}
	return out
}

func (a *AlmacenSQLC) ObtenerPrecio(id int) (*models.PrecioRecurso, error) {
	f, err := a.q.BuscarPrecioPorID(ctx(), int64(id))
	if err != nil {
		return nil, err
	}
	p := aPrecioDominio(f)
	return &p, nil
}

func (a *AlmacenSQLC) CrearPrecio(in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	f, err := a.q.CrearPrecio(ctx(), sqlcdb.CrearPrecioParams{
		RecursoTipo:   in.RecursoTipo,
		RecursoID:     int64(in.RecursoID),
		Precio:        in.Precio.String(),
		FechaVigencia: in.FechaVigencia.Format(time.RFC3339),
		CreatedAt:     formatTime(time.Now()),
	})
	if err != nil {
		return nil, err
	}
	p := aPrecioDominio(f)
	return &p, nil
}

func (a *AlmacenSQLC) ActualizarPrecio(id int, in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	f, err := a.q.ActualizarPrecio(ctx(), sqlcdb.ActualizarPrecioParams{
		RecursoTipo:   in.RecursoTipo,
		RecursoID:     int64(in.RecursoID),
		Precio:        in.Precio.String(),
		FechaVigencia: in.FechaVigencia.Format(time.RFC3339),
		ID:            int64(id),
	})
	if err != nil {
		return nil, err
	}
	p := aPrecioDominio(f)
	return &p, nil
}

func (a *AlmacenSQLC) EliminarPrecio(id int) error {
	filas, err := a.q.EliminarPrecio(ctx(), int64(id))
	if err != nil {
		return err
	}
	if filas == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (a *AlmacenSQLC) HistorialPrecios(tipo string, recursoID int) []*models.PrecioRecurso {
	filas, err := a.q.ListarPreciosPorRecurso(ctx(), sqlcdb.ListarPreciosPorRecursoParams{
		RecursoTipo: tipo,
		RecursoID:   int64(recursoID),
	})
	if err != nil {
		return nil
	}
	out := make([]*models.PrecioRecurso, 0, len(filas))
	for i := range filas {
		p := aPrecioDominio(filas[i])
		out = append(out, &p)
	}
	return out
}

func (a *AlmacenSQLC) PrecioVigente(tipo string, recursoID int) (*models.PrecioRecurso, error) {
	filas, err := a.q.ListarPreciosPorRecurso(ctx(), sqlcdb.ListarPreciosPorRecursoParams{
		RecursoTipo: tipo,
		RecursoID:   int64(recursoID),
	})
	if err != nil {
		return nil, err
	}
	if len(filas) == 0 {
		return nil, sql.ErrNoRows
	}
	p := aPrecioDominio(filas[0])
	return &p, nil
}

func (a *AlmacenSQLC) ExisteRecurso(tipo string, id int) error {
	switch tipo {
	case "material":
		_, err := a.q.BuscarMaterialPorID(ctx(), int64(id))
		return err
	case "mano_obra":
		_, err := a.q.BuscarManoObraPorID(ctx(), int64(id))
		return err
	case "equipo":
		_, err := a.q.BuscarEquipoPorID(ctx(), int64(id))
		return err
	default:
		return sql.ErrNoRows
	}
}

// =========================================================
// Usuario (UsuarioRepository)
// =========================================================

func (a *AlmacenSQLC) CrearUsuario(in models.EntradaUsuario) (*models.Usuario, error) {
	f, err := a.q.CrearUsuario(ctx(), sqlcdb.CrearUsuarioParams{
		Email:        in.Email,
		PasswordHash: in.Password,
		CreatedAt:    formatTime(time.Now()),
	})
	if err != nil {
		return nil, err
	}
	u := aUsuarioDominio(f)
	return &u, nil
}

func (a *AlmacenSQLC) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	u, err := a.q.BuscarUsuarioPorEmail(ctx(), email)
	if err != nil {
		return models.Usuario{}, false
	}
	return aUsuarioDominio(u), true
}

func (a *AlmacenSQLC) ListarUsuarios() []*models.Usuario {
	filas, err := a.q.ListarUsuarios(ctx())
	if err != nil {
		return nil
	}
	out := make([]*models.Usuario, 0, len(filas))
	for i := range filas {
		u := aUsuarioDominio(filas[i])
		out = append(out, &u)
	}
	return out
}

func (a *AlmacenSQLC) ObtenerUsuarioPorID(id int) (*models.Usuario, bool) {
	u, err := a.q.BuscarUsuarioPorID(ctx(), int64(id))
	if err != nil {
		return nil, false
	}
	v := aUsuarioDominio(u)
	return &v, true
}

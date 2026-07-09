package storage

import (
//<<<<<<< HEAD
	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"gorm.io/gorm"
)

// AlmacenSQLite implementa la interfaz Almacen usando GORM sobre SQLite.
//
// Fíjense: los métodos tienen EXACTAMENTE las mismas firmas que los de Memoria.
// Por eso el Server y los handlers no se enteran de cuál de los dos reciben.
type AlmacenSQLite struct {
	db *gorm.DB
}

// NuevoAlmacenSQLite envuelve una conexión *gorm.DB ya abierta.
func NuevoAlmacenSQLite(db *gorm.DB) *AlmacenSQLite {
	return &AlmacenSQLite{db: db}
}

// =========================================================
// INCIDENCIAS
// =========================================================

func (a *AlmacenSQLite) ListarIncidencias() []models.Incidencia {
	var incidencia []models.Incidencia
	a.db.Find(&incidencia)
	return incidencia
}

func (a *AlmacenSQLite) BuscarIncidenciaPorID(id int) (models.Incidencia, bool) {
	var i models.Incidencia
	if err := a.db.First(&i, id).Error; err != nil {
		// Absorbemos el error de la DB y conservamos la firma comma-ok.
		return models.Incidencia{}, false
	}
	return i, true
}

func (a *AlmacenSQLite) BuscarIncidenciaPorEntidad(id int, tipo string) (models.Incidencia, bool) {
	var i models.Incidencia
	if err := a.db.Where("entidad_tipo = ?", tipo).First(&i, id).Error; err != nil {
		// Absorbemos el error de la DB y conservamos la firma comma-ok.
		return models.Incidencia{}, false
	}
	return i, true
}

func (a *AlmacenSQLite) CrearIncidencia(i models.Incidencia) models.Incidencia {
	a.db.Create(&i) // GORM rellena el ID autogenerado en &p
	return i
}

func (a *AlmacenSQLite) ActualizarIncidencia(id int, datos models.Incidencia) (models.Incidencia, bool) {
	var existente models.Incidencia
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Incidencia{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarIncidencia(id int) bool {
	res := a.db.Delete(&models.Incidencia{}, id)
	return res.RowsAffected > 0
}

func (a *AlmacenSQLite) ListarObras() []models.Obra {
	var obras []models.Obra
	a.db.Find(&obras)
	return obras
}

func (a *AlmacenSQLite) BuscarObraPorID(id int) (models.Obra, bool) {
	var o models.Obra
	if err := a.db.First(&o, id).Error; err != nil {
		return models.Obra{}, false
	}
	return o, true
}

func (a *AlmacenSQLite) CrearObra(o models.Obra) models.Obra {
	a.db.Create(&o)
	return o
}

func (a *AlmacenSQLite) ActualizarObra(id int, datos models.Obra) (models.Obra, bool) {
	var existente models.Obra
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Obra{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarObra(id int) bool {
	res := a.db.Delete(&models.Obra{}, id)
	return res.RowsAffected > 0
}

func (a *AlmacenSQLite) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario

	if err := a.db.First(&u, email).Error; err != nil {
		return models.Usuario{}, false
	}

	return u, true
}

func (a *AlmacenSQLite) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	a.db.Create(&u)
	return u, nil
}

// Chequeo en tiempo de compilación: AlmacenSQLite debe cumplir Almacen.
var _ Almacen = (*AlmacenSQLite)(nil)


	"errors"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AlamcenSQlite struct {
	db *gorm.DB
}

func NewAlmacenSQLite(db *gorm.DB) *AlamcenSQlite {
	return &AlamcenSQlite{db: db}
}

func precioADecimal(precioStr string) (decimal.Decimal, error) {
	if precioStr == "" {
		return decimal.Zero, errors.New("precio vacio")
	}
	precio, err := decimal.NewFromString(precioStr)
	if err != nil {
		return decimal.Zero, err
	}
	return precio, nil
}

func (r *AlamcenSQlite) CrearMateriales(in models.EntradaMaterial) (*models.Material, error) {
	precio, err := precioADecimal(in.PrecioReferencia)
	if err != nil {
		return nil, err
	}
	mat := models.Material{
		Nombre:           in.Nombre,
		Descripcion:      in.Descripcion,
		Unidad:           in.Unidad,
		PrecioReferencia: precio,
	}
	if err := r.db.Create(&mat).Error; err != nil {
		return nil, err
	}
	return &mat, nil
}

func (r *AlamcenSQlite) ObtenerMateriales(id int) (*models.Material, bool) {
	var mat models.Material
	if err := r.db.First(&mat, id).Error; err != nil {
		return nil, false
	}
	return &mat, true
}

func (r *AlamcenSQlite) ListarMateriales() []*models.Material {
	var mats []models.Material
	r.db.Find(&mats)
	out := make([]*models.Material, len(mats))
	for i := range mats {
		out[i] = &mats[i]
	}
	return out
}

func (r *AlamcenSQlite) ActualizarMateriales(id int, in models.EntradaMaterial) (*models.Material, bool) {
	var mat models.Material
	if err := r.db.First(&mat, id).Error; err != nil {
		return nil, false
	}
	precio, err := precioADecimal(in.PrecioReferencia)
	if err != nil {
		return nil, false
	}
	mat.Nombre = in.Nombre
	mat.Descripcion = in.Descripcion
	mat.Unidad = in.Unidad
	mat.PrecioReferencia = precio
	if err := r.db.Save(&mat).Error; err != nil {
		return nil, false
	}
	return &mat, true
}

func (r *AlamcenSQlite) EliminarMateriales(id int) bool {
	res := r.db.Delete(&models.Material{}, id)
	if res.Error != nil || res.RowsAffected == 0 {
		return false
	}
	return true
}


//Mano Obra

func (r *AlamcenSQlite) CrearManoObra(in models.EntradaManoObra) (*models.ManoObra, error) {
	mo := models.ManoObra{
		Descripcion:     in.Descripcion,
		Categoria:       in.Categoria,
		Unidad:          in.Unidad,
		CostoReferencia: in.CostoReferencia,
	}
	if err := r.db.Create(&mo).Error; err != nil {
		return nil, err
	}
	return &mo, nil
}

func (r *AlamcenSQlite) ObtenerManoObra(id int) (*models.ManoObra, bool) {
	var mo models.ManoObra
	if err := r.db.First(&mo, id).Error; err != nil {
		return nil, false
	}
	return &mo, true
}

func (r *AlamcenSQlite) ListarManoObra() []*models.ManoObra {
	var items []models.ManoObra
	r.db.Find(&items)
	out := make([]*models.ManoObra, len(items))
	for i := range items {
		out[i] = &items[i]
	}
	return out
}

func (r *AlamcenSQlite) ActualizarManoObra(id int, in models.EntradaManoObra) (*models.ManoObra, bool) {
	var mo models.ManoObra
	if err := r.db.First(&mo, id).Error; err != nil {
		return nil, false
	}
	mo.Descripcion = in.Descripcion
	mo.Categoria = in.Categoria
	mo.Unidad = in.Unidad
	mo.CostoReferencia = in.CostoReferencia
	if err := r.db.Save(&mo).Error; err != nil {
		return nil, false
	}
	return &mo, true
}

func (r *AlamcenSQlite) EliminarManoObra(id int) bool {
	res := r.db.Delete(&models.ManoObra{}, id)
	if res.Error != nil || res.RowsAffected == 0 {
		return false
	}
	return true
}

//equipos


func (r *AlamcenSQlite) CrearEquipo(in models.EntradaEquipo) (*models.Equipo, error) {
	eq := models.Equipo{
		Nombre:     in.Nombre,
		Tipo:       in.Tipo,
		Unidad:     in.Unidad,
		CostoHora:  in.CostoHora,
		Disponible: in.Disponible,
	}
	if err := r.db.Create(&eq).Error; err != nil {
		return nil, err
	}
	return &eq, nil
}

func (r *AlamcenSQlite) ObtenerEquipo(id int) (*models.Equipo, error) {
	var eq models.Equipo
	if err := r.db.First(&eq, id).Error; err != nil {
		return nil, err
	}
	return &eq, nil
}

func (r *AlamcenSQlite) ListarEquipos() []*models.Equipo {
	var items []models.Equipo
	r.db.Find(&items)
	out := make([]*models.Equipo, len(items))
	for i := range items {
		out[i] = &items[i]
	}
	return out
}

func (r *AlamcenSQlite) ActualizarEquipo(id int, in models.EntradaEquipo) (*models.Equipo, error) {
	var eq models.Equipo
	if err := r.db.First(&eq, id).Error; err != nil {
		return nil, err
	}
	eq.Nombre = in.Nombre
	eq.Tipo = in.Tipo
	eq.Unidad = in.Unidad
	eq.CostoHora = in.CostoHora
	eq.Disponible = in.Disponible
	if err := r.db.Save(&eq).Error; err != nil {
		return nil, err
	}
	return &eq, nil
}

func (r *AlamcenSQlite) EliminarEquipo(id int) error {
	res := r.db.Delete(&models.Equipo{}, id)
	if res.Error != nil || res.RowsAffected == 0 {
		return res.Error
	}
	return nil
}


//Precios
func (r *AlamcenSQlite) ListarPrecios() []*models.PrecioRecurso {
	var items []models.PrecioRecurso
	r.db.Find(&items)
	out := make([]*models.PrecioRecurso, len(items))
	for i := range items {
		out[i] = &items[i]
	}
	return out
}

func (r *AlamcenSQlite) ObtenerPrecio(id int) (*models.PrecioRecurso, error) {
	var p models.PrecioRecurso
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *AlamcenSQlite) CrearPrecio(in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	p := models.PrecioRecurso{
		RecursoTipo:   in.RecursoTipo,
		RecursoID:     in.RecursoID,
		Precio:        in.Precio,
		FechaVigencia: in.FechaVigencia,
	}
	if err := r.db.Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *AlamcenSQlite) ActualizarPrecio(id int, in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	var p models.PrecioRecurso
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	p.Precio = in.Precio
	if err := r.db.Save(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *AlamcenSQlite) EliminarPrecio(id int) error {
	res := r.db.Delete(&models.PrecioRecurso{}, id)
	if res.Error != nil || res.RowsAffected == 0 {
		return res.Error
	}
	return nil
}

func (r *AlamcenSQlite) HistorialPrecios(tipo string, recursoID int) []*models.PrecioRecurso {
	var items []models.PrecioRecurso
	r.db.Where("recurso_tipo = ? AND recurso_id = ?", tipo, recursoID).Find(&items)
	out := make([]*models.PrecioRecurso, len(items))
	for i := range items {
		out[i] = &items[i]
	}
	return out
}

func (r *AlamcenSQlite) PrecioVigente(tipo string, recursoID int) (*models.PrecioRecurso, error) {
	var p models.PrecioRecurso
	if err := r.db.Where("recurso_tipo = ? AND recurso_id = ?", tipo, recursoID).Order("fecha_vigencia desc").First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *AlamcenSQlite) ExisteRecurso(tipo string, id int) error {
	switch tipo {
	case "material":
		var m models.Material
		if err := r.db.First(&m, id).Error; err != nil {
			return err
		}
	case "mano_obra":
		var m models.ManoObra
		if err := r.db.First(&m, id).Error; err != nil {
			return err
		}
	case "equipo":
		var e models.Equipo
		if err := r.db.First(&e, id).Error; err != nil {
			return err
		}
	default:
		return gorm.ErrRecordNotFound
	}
	return nil
}


// =========================================================
// SEEDS
// =========================================================
func (r *AlamcenSQlite) SembrarSiVacio() {
	var n int64
	r.db.Model(&models.Material{}).Count(&n)
	if n > 0 {
		return
	}

	now := time.Now()
	hace30 := now.AddDate(0, -1, 0)
	hace60 := now.AddDate(0, -2, 0)

	materiales := []struct{ nombre, descripcion, unidad, precio string }{
		{"Cemento Portland tipo I", "Saco de 50 kg para hormigón estructural", "unidad", "9.50"},
		{"Arena fina lavada", "Arena de río tamizada, libre de arcilla", "m³", "22.00"},
		{"Grava triturada 3/4\"", "Agregado grueso para hormigón", "m³", "28.50"},
		{"Varilla corrugada 12mm", "Acero de refuerzo ASTM A615 Gr.60", "kg", "1.15"},
		{"Varilla corrugada 16mm", "Acero de refuerzo ASTM A615 Gr.60", "kg", "1.12"},
		{"Bloque de hormigón 15x20x40", "Bloque vibrado para mampostería", "unidad", "0.65"},
		{"Ladrillo mambrón", "Ladrillo artesanal para fachada", "unidad", "0.28"},
		{"Tubo PVC presión 110mm", "Tubería PVC para agua potable, 6m", "unidad", "18.90"},
		{"Alambre de amarre #18", "Rollo de 25 kg", "kg", "1.80"},
		{"Pintura látex interior", "Pintura lavable blanca, 4 litros", "unidad", "12.80"},
		{"Cerámica piso 40x40", "Cerámica esmaltada para interior", "m²", "8.50"},
		{"Porcelanato 60x60", "Porcelanato rectificado mate", "m²", "18.00"},
		{"Impermeabilizante líquido", "Cementoso bicomponente, 20 kg", "unidad", "35.00"},
		{"Malla electrosoldada 15x15", "Acero para losas, rollo 6x2.4m", "unidad", "38.00"},
		{"Estuco listo para uso", "Pasta para empaste interior, 25 kg", "unidad", "7.20"},
	}
	for _, m := range materiales {
		r.db.Create(&models.Material{
			Nombre: m.nombre, Descripcion: m.descripcion,
			Unidad: m.unidad, PrecioReferencia: decimal.RequireFromString(m.precio),
		})

	}
	manoObras := []struct{ descripcion, categoria, unidad, costo string }{
		{"Maestro de obra general", "oficial", "día", "35.00"},
		{"Albañil - mampostería y enlucido", "oficial", "día", "28.00"},
		{"Fierrero - armado de acero", "oficial", "día", "30.00"},
		{"Carpintero de encofrado", "oficial", "día", "28.00"},
		{"Plomero instalaciones sanitarias", "oficial", "día", "32.00"},
		{"Electricista instalaciones", "oficial", "día", "32.00"},
		{"Pintor de interiores y exteriores", "oficial", "día", "26.00"},
		{"Ayudante de albañilería", "ayudante", "día", "18.00"},
		{"Ayudante de fierrero", "ayudante", "día", "18.00"},
		{"Peón de obra general", "ayudante", "día", "15.00"},
		{"Soldador estructura metálica", "especialista", "hora", "14.00"},
		{"Topógrafo de replanteo", "especialista", "día", "45.00"},
		{"Inspector de calidad hormigón", "especialista", "día", "50.00"},
	}
	for _, m := range manoObras {
		r.db.Create(&models.ManoObra{
			Descripcion:     m.descripcion,
			Categoria:       m.categoria,
			Unidad:          m.unidad,
			CostoReferencia: decimal.RequireFromString(m.costo),
		})
	}

	equipos := []struct {
		nombre, tipo, unidad, costoHora string
		disponible                      bool
	}{
		{"Concretera 1 saco eléctrica", "liviano", "hora", "8.50", true},
		{"Vibrador de hormigón 2\"", "liviano", "hora", "3.50", true},
		{"Amoladora angular 4.5\"", "liviano", "hora", "2.00", true},
		{"Compresor de aire 100 lt", "liviano", "hora", "4.00", true},
		{"Nivel láser rotativo", "liviano", "hora", "3.00", true},
		{"Excavadora sobre orugas CAT 320", "pesado", "hora", "85.00", true},
		{"Retroexcavadora JCB 3CX", "pesado", "hora", "65.00", true},
		{"Motoniveladora 140K", "pesado", "hora", "90.00", false},
		{"Compactadora tipo sapo", "pesado", "hora", "15.00", true},
		{"Volqueta 8 m³", "pesado", "hora", "55.00", true},
		{"Grúa torre 30m", "pesado", "hora", "120.00", false},
		{"Bomba de hormigón estacionaria", "pesado", "hora", "95.00", true},
	}
	for _, e := range equipos {
		r.db.Create(&models.Equipo{
			Nombre:        e.nombre,
			Tipo:          e.tipo,
			Unidad:        e.unidad,
			CostoHora:     decimal.RequireFromString(e.costoHora),
			Disponible:    e.disponible,
		})

	}

	type precioSeed struct {
		tipo   string
		id     int
		precio string
		fecha  time.Time
	}

	precios := []precioSeed{
		{"material", 1, "8.80", hace60},
		{"material", 1, "9.20", hace30},
		{"material", 1, "9.50", now},
		{"material", 2, "20.00", hace60},
		{"material", 2, "21.00", hace30},
		{"material", 2, "22.00", now},
		{"material", 4, "1.05", hace60},
		{"material", 4, "1.10", hace30},
		{"material", 4, "1.15", now},
		{"mano_obra", 16, "32.00", hace60},
		{"mano_obra", 16, "33.50", hace30},
		{"mano_obra", 16, "35.00", now},
		{"mano_obra", 17, "25.00", hace60},
		{"mano_obra", 17, "26.50", hace30},
		{"mano_obra", 17, "28.00", now},
		{"equipo", 34, "78.00", hace60},
		{"equipo", 34, "82.00", hace30},
		{"equipo", 34, "85.00", now},
		{"equipo", 35, "60.00", hace60},
		{"equipo", 35, "62.50", hace30},
		{"equipo", 35, "65.00", now},
	}
	for _, p := range precios {
		r.db.Create(&models.PrecioRecurso{
			RecursoTipo:   p.tipo,
			RecursoID:     p.id,
			Precio:        decimal.RequireFromString(p.precio),
			FechaVigencia: p.fecha,
		})

	}
}
//>>>>>> Modulo1/Catalogo

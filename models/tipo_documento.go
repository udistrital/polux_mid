package models

type TipoDocumento struct {
	Id                   int                   
	Nombre               string                
	Descripcion          string                
	CodigoAbreviacion    string                
	Activo               bool                  
	NumeroOrden          float64               
	Tamano               float64               
	Extension            string                
	Workspace            string                
	TipoDocumentoNuxeo   string                
	FechaCreacion        string                
	FechaModificacion    string                
	DominioTipoDocumento *DominioTipoDocumento 
}



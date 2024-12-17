package constant

const (
	SpanCreateProduct = "create_product"
	SpanGetProduct    = "get_product"
	SpanListProducts  = "list_products"
	SpanUpdateProduct = "update_product"
	SpanDeleteProduct = "delete_product"

	SpanCreateBrand = "create_brand"
	SpanGetBrand    = "get_brand"
	SpanListBrands  = "list_brands"
	SpanDeleteBrand = "delete_brand"

	AttributeProductID  = "product.id"
	AttributeBrandID    = "brand.id"
	AttributePageSize   = "page.size"
	AttributePageNumber = "page.number"
	AttributeUserAgent  = "http.user_agent"
	AttributeStatusCode = "http.status_code"
	AttributeMethod     = "http.method"
	AttributePath       = "http.path"

	MetricRequestDuration    = "http.request.duration"
	MetricRequestTotal       = "http.request.total"
	MetricDatabaseCallTotal  = "db.call.total"
	MetricDatabaseErrorTotal = "db.error.total"
	MetricProductCreated     = "product.created.total"
	MetricProductUpdated     = "product.updated.total"
	MetricProductDeleted     = "product.deleted.total"
	MetricBrandCreated       = "brand.created.total"
	MetricBrandDeleted       = "brand.deleted.total"

	DefaultServiceName    = "ecommerce-service"
	DefaultServiceVersion = "1.0.0"
	DefaultEnvironment    = "development"
	DefaultOTLPEndpoint   = "localhost:4317"
)

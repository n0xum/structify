export type TagEntry = {
  tag: string;
  description: string;
};

export type TagGroup = {
  group: string;
  tags: TagEntry[];
};

export type MethodEntry = {
  pattern: string;
  description: string;
};

export type MethodGroup = {
  group: string;
  methods: MethodEntry[];
};

/** Method naming conventions for the Interface Repository generator. */
export const METHOD_GROUPS: MethodGroup[] = [
  {
    group: "CRUD (any name)",
    methods: [
      { pattern: "Create(ctx, *T) (*T, error)", description: "INSERT … RETURNING" },
      { pattern: "GetByID(ctx, id) (*T, error)", description: "SELECT … WHERE id = $1" },
      { pattern: "Update(ctx, *T) error", description: "UPDATE … WHERE id = $n" },
      { pattern: "Delete(ctx, id) error", description: "DELETE … WHERE id = $1" },
      { pattern: "List(ctx) ([]*T, error)", description: "SELECT … ORDER BY id" },
    ],
  },
  {
    group: "FindBy (field conditions)",
    methods: [
      { pattern: "FindByField(ctx, v) (*T, error)", description: "WHERE field = $1 → single row" },
      { pattern: "FindByField(ctx, v) ([]*T, error)", description: "WHERE field = $1 → slice" },
      { pattern: "FindByXAndY(ctx, x, y) …", description: "WHERE x = $1 AND y = $2" },
      { pattern: "FindByXOrY(ctx, x, y) …", description: "WHERE x = $1 OR y = $2" },
    ],
  },
  {
    group: "SmartQuery (method name = full query)",
    methods: [
      { pattern: "ListXsByField(ctx, v) ([]*X, error)", description: "SELECT … WHERE field = $1" },
      { pattern: "CountXsByField(ctx, v) (int64, error)", description: "SELECT COUNT(*) WHERE …" },
      { pattern: "ExistsXByField(ctx, v) (bool, error)", description: "SELECT EXISTS(SELECT 1 …)" },
      { pattern: "DeleteXByField(ctx, v) error", description: "DELETE … WHERE field = $1" },
      { pattern: "…OrderByFieldDesc", description: "Append ORDER BY field DESC" },
      { pattern: "…OrderByFieldAsc", description: "Append ORDER BY field ASC" },
    ],
  },
  {
    group: "Operator suffixes",
    methods: [
      { pattern: "ByAgeGreaterThan", description: "age > $n" },
      { pattern: "ByAgeLessThan", description: "age < $n" },
      { pattern: "ByNameLike", description: "name LIKE $n" },
      { pattern: "ByIDIn", description: "id IN ($n)" },
      { pattern: "ByIDNotIn", description: "id NOT IN ($n)" },
      { pattern: "ByDeletedAtIsNull", description: "deleted_at IS NULL" },
      { pattern: "ByDeletedAtIsNotNull", description: "deleted_at IS NOT NULL" },
    ],
  },
  {
    group: "CustomSQL",
    methods: [
      { pattern: `//sql:"SELECT …"`, description: "Exact SQL used verbatim; $1, $2 … for params" },
    ],
  },
];

export const TAG_GROUPS: TagGroup[] = [
  {
    group: "Basic",
    tags: [
      { tag: `db:"pk"`, description: "Primary key" },
      { tag: `db:"unique"`, description: "UNIQUE constraint" },
      { tag: `db:"-"`, description: "Exclude field" },
      { tag: `db:"table:name"`, description: "Override table name (on struct)" },
    ],
  },
  {
    group: "Constraints",
    tags: [
      { tag: `db:"check:age >= 18"`, description: "CHECK constraint" },
      { tag: `db:"default:true"`, description: "DEFAULT value" },
      { tag: `db:"enum:a,b,c"`, description: "Enum CHECK (IN clause)" },
    ],
  },
  {
    group: "Indexes",
    tags: [
      { tag: `db:"index"`, description: "Auto-named index" },
      { tag: `db:"index:idx_name"`, description: "Named index" },
      { tag: `db:"unique_index"`, description: "Unique index" },
      { tag: `db:"unique_index:uq_name"`, description: "Named unique index" },
    ],
  },
  {
    group: "Foreign Keys",
    tags: [
      { tag: `db:"fk:users,id"`, description: "FK → users(id)" },
      { tag: `db:"fk:users,id,on_delete:CASCADE"`, description: "FK with ON DELETE CASCADE" },
      { tag: `db:"fk:name,table,col"`, description: "Composite FK (same name groups cols)" },
    ],
  },
  {
    group: "Composite",
    tags: [
      { tag: `db:"pk" (multiple fields)`, description: "Composite primary key" },
      { tag: `db:"unique:uq_name"`, description: "Composite unique (same name groups cols)" },
    ],
  },
];

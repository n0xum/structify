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
    group: "CRUD methods",
    methods: [
      { pattern: "Create(ctx, *T) error", description: "Insert one row from struct fields." },
      { pattern: "Create(ctx, *T) (*T, error)", description: "Insert and return the saved row." },
      { pattern: "FindByID(ctx, id) (*T, error)", description: "Lookup by primary key." },
      { pattern: "Update(ctx, *T) error", description: "Update row matched by primary key." },
      { pattern: "DeleteByID(ctx, id) error", description: "Delete row matched by primary key." },
    ],
  },
  {
    group: "FindBy filters",
    methods: [
      { pattern: "FindByEmail(ctx, email) (*T, error)", description: "Single-row result by one field." },
      { pattern: "FindByStatus(ctx, status) ([]*T, error)", description: "Multi-row result by one field." },
      { pattern: "FindByEmailAndStatus(ctx, email, status) …", description: "Combine conditions with AND." },
      { pattern: "FindByEmailOrUsername(ctx, email, username) …", description: "Combine conditions with OR." },
    ],
  },
  {
    group: "SmartQuery methods",
    methods: [
      { pattern: "ListUsersByStatus(ctx, status) ([]*User, error)", description: "List with filter conditions." },
      { pattern: "CountUsersByStatus(ctx, status) (int64, error)", description: "Count rows matching filters." },
      { pattern: "ExistsUserByEmail(ctx, email) (bool, error)", description: "Return true when a row exists." },
      { pattern: "DeleteUserByEmail(ctx, email) error", description: "Delete rows matched by filters." },
      { pattern: "...OrderByCreatedDesc", description: "Add ORDER BY created DESC." },
      { pattern: "...OrderByCreatedAsc", description: "Add ORDER BY created ASC." },
    ],
  },
  {
    group: "Operator suffixes",
    methods: [
      { pattern: "ByAgeGreaterThan", description: "Use > comparison." },
      { pattern: "ByAgeLessThan", description: "Use < comparison." },
      { pattern: "ByNameLike", description: "Use LIKE comparison." },
      { pattern: "ByIDIn", description: "Use IN (...) comparison." },
      { pattern: "ByIDNotIn", description: "Use NOT IN (...) comparison." },
      { pattern: "ByDeletedAtIsNull", description: "Use IS NULL check." },
      { pattern: "ByDeletedAtIsNotNull", description: "Use IS NOT NULL check." },
    ],
  },
  {
    group: "Custom SQL override",
    methods: [
      { pattern: `//sql:"SELECT ... WHERE email = $1"`, description: "Run explicit SQL when naming rules are not enough." },
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

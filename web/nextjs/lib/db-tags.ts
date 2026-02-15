export type TagEntry = {
  tag: string;
  description: string;
};

export type TagGroup = {
  group: string;
  tags: TagEntry[];
};

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

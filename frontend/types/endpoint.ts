export interface SchemaField {
  name: string;
  type: string;
}

export interface Endpoint {
  id: string; // Or number, depending on backend (assuming string/uuid or number)
  path: string;
  method: string;
  delay?: number;
  headers?: Record<string, string>;
  count: number;
  schema: SchemaField[];
  createdAt?: string;
  updatedAt?: string;
}

export interface EndpointPayload {
  path: string;
  method: string;
  delay?: number;
  headers?: Record<string, string>;
  count: number;
  schema: SchemaField[];
}

export interface MockRecord {
  id?: string | number;
  [key: string]: any;
}

export interface EndpointDetail extends Endpoint {
  records: MockRecord[];
}

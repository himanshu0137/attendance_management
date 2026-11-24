export interface Employee {
  id: string;
  keycloakSub: string;
  role: string;
}

export interface CreateEmployeeRequest {
  tenantId: string;
  firstName: string;
  lastName: string;
  email: string;
  employeeCode: string;
  designation: string;
  department: string;
  managerId?: string;
  branchId?: string;
}

export interface Document {
  name: string;
  type: string;
  description: string;
  url: string;
}

export interface EmployeeWithDocuments {
  id: string;
  keycloakSub: string;
  role: string;
  documents: Document[];
  joiningDate: Date;
  position: string;
  department: string;
}

export interface OnboardEmployeeReq {
  documents: Document[];
  joiningDate: Date;
  position: string;
  department: string;
}
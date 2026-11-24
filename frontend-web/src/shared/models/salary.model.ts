export interface CreateSalaryStructureRequest {
  payBasis: string; // "monthly" or "daily"
  basic: number;
  hra: number;
  da: number;
  conveyance: number;
  specialAllowance: number;
  medicalAllowance: number;
  bonus: number;
  overtimeRate: number;
  effectiveDate: Date;
  createdAt: Date;
}

export interface SalaryStructure {
  id: string;
  employeeId: string;
  tenantId: string;
  payBasis: string; // "monthly" or "daily"
  basic: number;
  hra: number;
  da: number;
  conveyance: number;
  specialAllowance: number;
  medicalAllowance: number;
  bonus: number;
  overtimeRate: number;
  effectiveDate: Date;
  createdAt: Date;
  updatedAt: Date;
}
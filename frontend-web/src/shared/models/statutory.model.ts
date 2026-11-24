export interface StatutoryCompliance {
  id: string;
  tenantId: string;
  complianceType: string; // e.g., "ESIC", "EPF", "Professional_Tax", "Income_Tax"
  applicableFrom: Date;
  applicableTo?: Date;
  rate: number;
  maxLimit?: number;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
}

export interface EmployeeStatutoryDeduction {
  id: string;
  employeeId: string;
  statutoryId: string;
  applicableFrom: Date;
  applicableTo?: Date;
  rate: number;
  maxLimit?: number;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
}
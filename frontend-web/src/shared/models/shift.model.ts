export interface Shift {
  id: string;
  tenantId: string;
  name: string;
  description: string;
  startTime: Date;
  endTime: Date;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
}
export interface AttendanceRule {
  id: string;
  tenantId: string;
  name: string;
  description: string;
  ruleType: string; // e.g., "late_arrival", "early_departure", "overtime"
  condition: Record<string, any>;
  action: Record<string, any>;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
}
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { OrganizationService } from '../../services/organization.service';
import { CreateSalaryStructureRequest } from '../../shared/models/salary.model';

@Component({
  selector: 'app-salary-structure',
  templateUrl: './salary-structure.component.html',
  styleUrls: ['./salary-structure.component.scss']
})
export class SalaryStructureComponent implements OnInit {
  salaryForm: FormGroup;
  isSubmitting = false;
  tenantId: string = '';
  employeeId: string = '';

  constructor(
    private fb: FormBuilder,
    private organizationService: OrganizationService
  ) {
    this.salaryForm = this.fb.group({
      tenantId: ['', Validators.required],
      employeeId: ['', Validators.required],
      payBasis: ['monthly', Validators.required],
      basic: [0, Validators.required],
      hra: [0],
      da: [0],
      conveyance: [0],
      specialAllowance: [0],
      medicalAllowance: [0],
      bonus: [0],
      overtimeRate: [0],
      effectiveDate: ['', Validators.required]
    });
  }

  ngOnInit(): void {}

  onSubmit(): void {
    if (this.salaryForm.valid) {
      this.isSubmitting = true;
      const salaryData: CreateSalaryStructureRequest = this.salaryForm.value;
      
      // Convert date string to Date object if needed
      if (typeof salaryData.effectiveDate === 'string') {
        salaryData.effectiveDate = new Date(salaryData.effectiveDate);
      }
      
      this.organizationService.createSalaryStructure(
        salaryData.tenantId, 
        salaryData.employeeId, 
        salaryData
      ).subscribe({
        next: (salaryStructure) => {
          console.log('Salary structure created:', salaryStructure);
          alert('Salary structure created successfully!');
          this.salaryForm.reset();
          this.isSubmitting = false;
        },
        error: (error) => {
          console.error('Error creating salary structure:', error);
          alert('Error creating salary structure: ' + error);
          this.isSubmitting = false;
        }
      });
    }
  }

  getSalaryStructure(): void {
    if (this.tenantId && this.employeeId) {
      this.organizationService.getSalaryStructure(this.tenantId, this.employeeId).subscribe({
        next: (salaryStructure) => {
          console.log('Retrieved salary structure:', salaryStructure);
          // Populate form with retrieved data
          this.salaryForm.patchValue({
            tenantId: this.tenantId,
            employeeId: this.employeeId,
            payBasis: salaryStructure.payBasis,
            basic: salaryStructure.basic,
            hra: salaryStructure.hra,
            da: salaryStructure.da,
            conveyance: salaryStructure.conveyance,
            specialAllowance: salaryStructure.specialAllowance,
            medicalAllowance: salaryStructure.medicalAllowance,
            bonus: salaryStructure.bonus,
            overtimeRate: salaryStructure.overtimeRate,
            effectiveDate: salaryStructure.effectiveDate
          });
        },
        error: (error) => {
          console.error('Error retrieving salary structure:', error);
          alert('Error retrieving salary structure: ' + error);
        }
      });
    }
  }

  updateSalaryStructure(): void {
    if (this.tenantId && this.employeeId && this.salaryForm.valid) {
      this.isSubmitting = true;
      const salaryData: CreateSalaryStructureRequest = this.salaryForm.value;
      
      // Convert date string to Date object if needed
      if (typeof salaryData.effectiveDate === 'string') {
        salaryData.effectiveDate = new Date(salaryData.effectiveDate);
      }
      
      this.organizationService.updateSalaryStructure(
        salaryData.tenantId, 
        salaryData.employeeId, 
        salaryData
      ).subscribe({
        next: (salaryStructure) => {
          console.log('Salary structure updated:', salaryStructure);
          alert('Salary structure updated successfully!');
          this.isSubmitting = false;
        },
        error: (error) => {
          console.error('Error updating salary structure:', error);
          alert('Error updating salary structure: ' + error);
          this.isSubmitting = false;
        }
      });
    }
  }
}
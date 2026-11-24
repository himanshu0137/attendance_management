import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { OrganizationService } from '../../services/organization.service';
import { CreateEmployeeRequest } from '../../shared/models/employee.model';

@Component({
  selector: 'app-employee-form',
  templateUrl: './employee-form.component.html',
  styleUrls: ['./employee-form.component.scss']
})
export class EmployeeFormComponent implements OnInit {
  employeeForm: FormGroup;
  isSubmitting = false;

  constructor(
    private fb: FormBuilder,
    private organizationService: OrganizationService
  ) {
    this.employeeForm = this.fb.group({
      tenantId: ['', Validators.required],
      firstName: ['', Validators.required],
      lastName: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      employeeCode: ['', Validators.required],
      designation: ['', Validators.required],
      department: ['', Validators.required],
      managerId: [''],
      branchId: ['']
    });
  }

  ngOnInit(): void {}

  onSubmit(): void {
    if (this.employeeForm.valid) {
      this.isSubmitting = true;
      const employeeData: CreateEmployeeRequest = this.employeeForm.value;
      
      this.organizationService.createEmployee(employeeData.tenantId, employeeData).subscribe({
        next: (employee) => {
          console.log('Employee created:', employee);
          alert('Employee created successfully!');
          this.employeeForm.reset();
          this.isSubmitting = false;
        },
        error: (error) => {
          console.error('Error creating employee:', error);
          alert('Error creating employee: ' + error);
          this.isSubmitting = false;
        }
      });
    }
  }
}
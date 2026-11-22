import { Component, Input } from '@angular/core';
import {
  MatDatepicker,
  MatDatepickerModule,
  MatFormField,
  MatFormFieldModule,
  MatInput,
  MatInputModule,
} from '@angular/material';

// DateRange is the pair of dates accepted by the reports endpoint.
export interface DateRange {
  date_from: string;
  date_to: string;
}

// DateRangePicker is a Material date-range picker bound to a DateRange.
@Component({
  selector: 'app-date-range',
  standalone: true,
  imports: [MatFormFieldModule, MatInputModule, MatDatepickerModule],
  templateUrl: './date-range.html',
  styleUrl: './date-range.scss'
})
export class DateRangePicker {
  @Input() value!: DateRange;
  @Input() disabled = false;

  start = new Date();
  end = new Date();

  apply() {
    const s = new Date(this.start.getFullYear(), this.start.getMonth(), this.start.getDate());
    const e = new Date(this.end.getFullYear(), this.end.getMonth(), this.end.getDate());
    this.value = {
      date_from: s.toISOString().slice(0, 10),
      date_to: e.toISOString().slice(0, 10),
    };
  }
}

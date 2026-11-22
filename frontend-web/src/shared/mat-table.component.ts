import { Component, Input, Output, EventEmitter } from '@angular/core';
import {
  MatColumnDef,
  MatHeaderCell,
  MatSort,
  MatSortChange,
  MatTable,
  MatTableModule,
} from '@angular/material';

// TableColumn describes one column of the shared table.
export interface TableColumn {
  key: string;
  label: string;
  sortable?: boolean;
}

// MatTable renders rows of objects with a column definition. Sorting is
// applied client-side and paging is optional.
@Component({
  selector: 'app-table',
  standalone: true,
  imports: [MatTableModule],
  templateUrl: './mat-table.html',
  styleUrl: './mat-table.scss'
})
export class MatTableComponent {
  @Input() columns: TableColumn[] = [];
  @Input() rows: object[] = [];
  @Input() page = 1;
  @Input() pageSize = 25;
  @Input() pageSizeOptions = [25, 50, 100];

  sortedRows: object[] = [];
  sortKey = '';
  sortDir = 'asc';

  @Output() pageChanged = new EventEmitter<number>();

  sort(key: string) {
    if (this.sortKey === key) {
      this.sortDir = this.sortDir === 'asc' ? 'desc' : 'asc';
    } else {
      this.sortKey = key;
      this.sortDir = 'asc';
    }
    const dir = this.sortDir === 'asc' ? 1 : -1;
    this.sortedRows = [...this.rows].sort((a, b) => {
      const av = a[this.sortKey];
      const bv = b[this.sortKey];
      if (av === bv) return 0;
      if (av === null || av === undefined) return 1;
      if (bv === null || bv === undefined) return -1;
      return (typeof av === 'number' ? av : String(av)).localeCompare(
        typeof bv === 'number' ? bv : String(bv),
        undefined,
        dir > 0 ? undefined : 'numeric'
      ) * dir;
    });
  }

  displayedRows = this.sortedRows;
}

import { Component, Input, signal } from '@angular/core';
import { MatBadge, MatBadgeModule } from '@angular/material';

// Colour applied to the badge face.
export type BadgeTone = 'primary' | 'accent' | 'warn' | 'error' | 'success' | 'neutral';

// StatusBadge renders a Material badge whose tone is driven by a signal.
@Component({
  selector: 'app-status-badge',
  standalone: true,
  imports: [MatBadgeModule],
  templateUrl: './badge.html',
  styleUrl: './badge.scss'
})
export class StatusBadge {
  @Input() label!: string;
  @Input() tone = 'neutral';
  badgeTone = signal<BadgeTone>('neutral');
}

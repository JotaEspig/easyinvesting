import { Component, OnInit } from '@angular/core';
import { HelloService } from '../services/hello.service';

@Component({
  selector: 'app-hello',
  imports: [],
  templateUrl: './hello.component.html',
  styleUrl: './hello.component.css'
})
export class HelloComponent {
  message: string = '';

  constructor(private helloService: HelloService) {}

  ngOnInit(): void {
    this.helloService.getHello().subscribe({
      next: (data) => this.message = data.message,
      error: (err) => console.error('API error:', err)
    });
}
}

import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HelloauthComponent } from './helloauth.component';

describe('HelloauthComponent', () => {
  let component: HelloauthComponent;
  let fixture: ComponentFixture<HelloauthComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HelloauthComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(HelloauthComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

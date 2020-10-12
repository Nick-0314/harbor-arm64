import { TestBed, inject } from '@angular/core/testing';

import { AccessLogService, AccessLogDefaultService } from './access-log.service';
import { SharedModule } from '../shared/shared.module';
import { SERVICE_CONFIG, IServiceConfig } from '../service.config';

describe('AccessLogService', () => {
  const mockConfig: IServiceConfig = {
    logBaseEndpoint: "/api/logs/testing"
  };

  let config: IServiceConfig;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        SharedModule
      ],
      providers: [
        AccessLogDefaultService,
        {
          provide: AccessLogService,
          useClass: AccessLogDefaultService
        }, {
          provide: SERVICE_CONFIG,
          useValue: mockConfig
        }]
    });

    config = TestBed.get(SERVICE_CONFIG);
  });

  it('should be initialized', inject([AccessLogDefaultService], (service: AccessLogService) => {
    expect(service).toBeTruthy();
  }));

  it('should inject the right config', () => {
    expect(config).toBeTruthy();
    expect(config.logBaseEndpoint).toEqual("/api/logs/testing");
  });

});

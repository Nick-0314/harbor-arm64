import { TestBed, inject } from '@angular/core/testing';

import { ConfigurationService, ConfigurationDefaultService } from './configuration.service';
import { SharedModule } from '../shared/shared.module';
import { SERVICE_CONFIG, IServiceConfig } from '../service.config';

describe('ConfigurationService', () => {
  const mockConfig: IServiceConfig = {
    configurationEndpoint: "/api/configurations/testing"
  };

  let config: IServiceConfig;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        SharedModule
      ],
      providers: [
        ConfigurationDefaultService,
        {
          provide: ConfigurationService,
          useClass: ConfigurationDefaultService
        }, {
          provide: SERVICE_CONFIG,
          useValue: mockConfig
        }]
    });

    config = TestBed.get(SERVICE_CONFIG);
  });

  it('should be initialized', inject([ConfigurationDefaultService], (service: ConfigurationService) => {
    expect(service).toBeTruthy();
  }));

  it('should inject the right config', () => {
    expect(config).toBeTruthy();
    expect(config.configurationEndpoint).toEqual("/api/configurations/testing");
  });

});

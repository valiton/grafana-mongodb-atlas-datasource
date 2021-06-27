import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { DataSourceOptions, SecureJsonData } from './types';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<DataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onPublicKeyChange = (event: any) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      publicKey: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  onPrivateKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        privateKey: event.target.value,
      },
    });
  };

  onResetPrivateKey = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        privateKey: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        privateKey: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const { jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as SecureJsonData;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="Public Key"
            onChange={this.onPublicKeyChange}
            value={jsonData.publicKey || ''}
            placeholder="json field returned to frontend"
          />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.privateKey) as boolean}
              value={secureJsonData.privateKey || ''}
              label="Private Key"
              placeholder="secure json field (backend only)"
              labelWidth={6}
              inputWidth={20}
              onReset={this.onResetPrivateKey}
              onChange={this.onPrivateKeyChange}
            />
          </div>
        </div>
      </div>
    );
  }
}

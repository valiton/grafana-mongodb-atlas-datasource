import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms, RadioButtonGroup, Field, Label } from '@grafana/ui';
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

  onApiTypeChange = (value: any) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      apiType: value,
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
      <>
        <div className="gf-form-group">
            <Label description={<span>Please enter your MongoDB programmatic API key as described <a href="https://docs.atlas.mongodb.com/configure-api-access" target="_blank">here</a></span>}>API Access Credentials</Label>
            <div className="gf-form">
            <FormField
              label="Public Key"
              onChange={this.onPublicKeyChange}
              value={jsonData.publicKey || ''}
              placeholder="e.g. wgfyfpdb"
            />
            </div>

            <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.privateKey) as boolean}
              value={secureJsonData.privateKey || ''}
              label="Private Key"
              labelWidth={6}
              inputWidth={20}
              onReset={this.onResetPrivateKey}
              onChange={this.onPrivateKeyChange}
            /></div>
        </div>

        <div className="gf-form-group">
          <div className="gf-form">
            <Field
              label="API type"
              description="Select whether you are using the free (public) version of the MongoDB cloud or the managed one (Atlas)"
            >
              <RadioButtonGroup
                options={[
                  { label: 'MongoDB Atlas Cloud API', value: 'atlas' },
                  { label: 'MongoDB Public Cloud API', value: 'public' },
                ]}
                value={jsonData.apiType || 'atlas'}
                onChange={this.onApiTypeChange}
              />
            </Field>
          </div>
        </div>
      </>
    );
  }
}

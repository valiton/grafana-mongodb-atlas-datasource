import { GenericDatasource, mapToTextValue } from "../src/datasource";

describe('GenericDatasource', function() {
  let ctx: any = {};

  beforeEach(function() {
    ctx.backendSrv = {};
    ctx.templateSrv = {
      replace: jest.fn().mockImplementation(value => value)
    };
    ctx.ds = new GenericDatasource({}, ctx.backendSrv, ctx.templateSrv);
  });

  describe('When invoking a query', () => {

    it('should return an empty array when no targets are set', done => {
      ctx.ds.query({targets: []}).then(result => {
        expect(result.data).toHaveLength(0);
        done();
      });
    });

    it('should return the server results when a target is set', done => {
      ctx.backendSrv.datasourceRequest = jest.fn().mockResolvedValue({
        data: { results: [{
          refId: 'A',
          series: [{
            name: 'X',
            points: [1, 2, 3]
          }]
        }]}
      });

      ctx.ds.query({ targets: ['hits'], range: { from: 'now-1h', to: 'now' }}).then(result => {
        var series = result.data[0];
        expect(series.target).toBe('X');
        expect(series.datapoints).toHaveLength(3);
        done();
      });
    });
  });

  describe('When invoking a metricFindQuery', () => {

    it ('should return the metric results when a target is null', done => {
      ctx.backendSrv.datasourceRequest = jest.fn().mockResolvedValue({
        data: [
          "metric_0",
          "metric_1",
          "metric_2",
        ]
      });

      ctx.ds.metricFindQuery({target: null}).then(function(result) {
        expect(result).toHaveLength(3);
        expect(result[0].text).toBe('metric_0');
        expect(result[0].value).toBe('metric_0');
        expect(result[1].text).toBe('metric_1');
        expect(result[1].value).toBe('metric_1');
        expect(result[2].text).toBe('metric_2');
        expect(result[2].value).toBe('metric_2');
        done();
      });
    });

    it ('should return the metric target results when a target is set', done => {
      ctx.backendSrv.datasourceRequest = jest.fn().mockImplementation(request => {
        var target = request.data.target;
        var result = [target + "_0", target + "_1", target + "_2"];

        return Promise.resolve({
          _request: request,
          data: result
        });
      });

      ctx.ds.metricFindQuery('search').then(function(result) {
        expect(result).toHaveLength(3);
        expect(result[0].text).toBe('search_0');
        expect(result[0].value).toBe('search_0');
        expect(result[1].text).toBe('search_1');
        expect(result[1].value).toBe('search_1');
        expect(result[2].text).toBe('search_2');
        expect(result[2].value).toBe('search_2');
        done();
      });
    });

    it ('should return the metric results when the target is an empty string', done => {
      ctx.backendSrv.datasourceRequest = jest.fn().mockResolvedValue({
        data: [
          "metric_0",
          "metric_1",
          "metric_2",
        ]
      });

      ctx.ds.metricFindQuery('').then(function(result) {
        expect(result).toHaveLength(3);
        expect(result[0].text).toBe('metric_0');
        expect(result[0].value).toBe('metric_0');
        expect(result[1].text).toBe('metric_1');
        expect(result[1].value).toBe('metric_1');
        expect(result[2].text).toBe('metric_2');
        expect(result[2].value).toBe('metric_2');
        done();
      });
    });

    it ('should return the metric results when the args are an empty object', done => {
      ctx.backendSrv.datasourceRequest = jest.fn().mockResolvedValue({
        data: [
          "metric_0",
          "metric_1",
          "metric_2",
        ]
      });

      ctx.ds.metricFindQuery().then(function(result) {
        expect(result).toHaveLength(3);
        expect(result[0].text).toBe('metric_0');
        expect(result[0].value).toBe('metric_0');
        expect(result[1].text).toBe('metric_1');
        expect(result[1].value).toBe('metric_1');
        expect(result[2].text).toBe('metric_2');
        expect(result[2].value).toBe('metric_2');
        done();
      });
    });

    it ('should return the metric target results when the args are a string', done => {
      ctx.backendSrv.datasourceRequest = jest.fn().mockImplementation(request => {
        var target = request.data.target;
        var result = [target + "_0", target + "_1", target + "_2"];

        return Promise.resolve({
          _request: request,
          data: result
        });
      });

      ctx.ds.metricFindQuery('search').then(function(result) {
        expect(result).toHaveLength(3);
        expect(result[0].text).toBe('search_0');
        expect(result[0].value).toBe('search_0');
        expect(result[1].text).toBe('search_1');
        expect(result[1].value).toBe('search_1');
        expect(result[2].text).toBe('search_2');
        expect(result[2].value).toBe('search_2');
        done();
      });
    });
  });

  describe('When mapping result to text value', () => {

    it ('should return data as text and as value', done => {
      var result = mapToTextValue({data: ["zero", "one", "two"]});

      expect(result).toHaveLength(3);
      expect(result[0].text).toBe('zero');
      expect(result[0].value).toBe('zero');
      expect(result[1].text).toBe('one');
      expect(result[1].value).toBe('one');
      expect(result[2].text).toBe('two');
      expect(result[2].value).toBe('two');
      done();
    });

    it ('should return text as text and value as value', done => {
      var data = [
        {text: "zero", value: "value_0"},
        {text: "one", value: "value_1"},
        {text: "two", value: "value_2"},
      ];

      var result = mapToTextValue({data: data});

      expect(result).toHaveLength(3);
      expect(result[0].text).toBe('zero');
      expect(result[0].value).toBe('value_0');
      expect(result[1].text).toBe('one');
      expect(result[1].value).toBe('value_1');
      expect(result[2].text).toBe('two');
      expect(result[2].value).toBe('value_2');
      done();
    });

    it ('should return data as text and index as value', done => {
      var data = [
        {a: "zero", b: "value_0"},
        {a: "one", b: "value_1"},
        {a: "two", b: "value_2"},
      ];

      var result = mapToTextValue({data: data});

      expect(result).toHaveLength(3);
      expect(result[0].text).toBe(data[0]);
      expect(result[0].value).toBe(0);
      expect(result[1].text).toBe(data[1]);
      expect(result[1].value).toBe(1);
      expect(result[2].text).toBe(data[2]);
      expect(result[2].value).toBe(2);
      done();
    });
  });
});

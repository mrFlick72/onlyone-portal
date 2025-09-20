import express, {Express} from 'express';
import dotenv from 'dotenv';
import { registerHealthEndPointFor } from './infrastructure/health';
import oAuth2ResourceServerMiddlewareFactory from "./infrastructure/oauth2Middleware";

import expressContext from 'express-request-context';
import {registerAccountEndPointFor} from "./api/AccountEndPoint";
var cors = require('cors')

process.on('uncaughtException', function (error) {
    console.log(error.stack); 
});
var corsOptions = {
    credentials: true,
    origin: [process.env.ALLOWED_ORIGIN],
}
dotenv.config();
const app: Express = express();

const promMid = require('express-prometheus-middleware');
app.use(promMid({
  metricsPath: '/metrics',
  collectDefaultMetrics: true,
  requestDurationBuckets: [0.1, 0.5, 1, 1.5],
  requestLengthBuckets: [512, 1024, 5120, 10240, 51200, 102400],
  responseLengthBuckets: [512, 1024, 5120, 10240, 51200, 102400],
  /**
   * Uncomenting the `authenticate` callback will make the `metricsPath` route
   * require authentication. This authentication callback can make a simple
   * basic auth test, or even query a remote server to validate access.
   * To access /metrics you could do:
   * curl -X GET user:password@localhost:9091/metrics
   */
  // authenticate: req => req.headers.authorization === 'Basic dXNlcjpwYXNzd29yZA==',
  /**
   * Uncommenting the `extraMasks` config will use the list of regexes to
   * reformat URL path names and replace the values found with a placeholder value
  */
  // extraMasks: [/..:..:..:..:..:../],
  /**
   * The prefix option will cause all metrics to have the given prefix.
   * E.g.: `app_prefix_http_requests_total`
   */
  // prefix: 'app_prefix_',
  /**
   * Can add custom labels with customLabels and transformLabels options
   */
  // customLabels: ['contentType'],
  // transformLabels(labels, req) {
  //   // eslint-disable-next-line no-param-reassign
  //   labels.contentType = req.headers['content-type'];
  // },
}));

let port = Number(process.env.APPLICATION_PORT) || 3000;

const actuatorApp: Express = express();
const actuatorPort = Number(process.env.ACTUATOR_PORT) || 3001;
app.use(expressContext())
app.use(cors(corsOptions));
app.use(express.json());
app.use(oAuth2ResourceServerMiddlewareFactory(process.env.ISSUER || ""));

registerHealthEndPointFor(actuatorApp)
registerAccountEndPointFor(app)

export {app, actuatorApp, port, actuatorPort}
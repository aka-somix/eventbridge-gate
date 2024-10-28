
/* eslint-disable @typescript-eslint/naming-convention */
/*
 * CONFIGURATION FILE
 * --------------------------------------------------------------------
 * This File holds the configuration parameters for this service,
 * providing a single source for fine tuning the service over its life span.
 *
 * ---
 * GUIDELINES
 * .
 ? How I decide wether to put or not a parameter in this file?
 *
 * When configuring an external source, ask yourself what are the parameters you
 * may want to change in the future, or experiment by passing them different values
 * from Env variables? Those are the pieces you want to include here.
 * .
 ! Do not include in this file any secret or sensible data. (Use the secret manager util)
 ! Do not fragment your configuration.
 * Avoid putting for the same service some parameters here and some parameters in place where
 * the service is initialized.
 *
 * ---
 ? How should I write a new configuration parameter?
 *
 * Every one of the parameters here below SHOULD BE:
 *
 * - Formatted in Upper Snake Case (THIS_IS_AN_EXAMPLE)
 * - Overridable from Environment Variable, yet have a default value (if it does makes sense to)
 * - Exported singularly. If a group of variable are always exposed together, consider to use objects
 */

import { LogLevel } from "@aws-lambda-powertools/logger/lib/cjs/types/Logger";

export const ENV = process.env.ENV;

export const AWS_API_VERSION = process.env.AWS_API_VERSION ?? "2017-10-17";

export const AWS_REGION = process.env.AWS_REGION ?? "eu-west-1";

export const LOG_LEVEL = (process.env.LOG_LEVEL ?? "INFO") as LogLevel;

// Base path of this service
export const BASE_PATH = process.env.BASE_PATH ?? "/main";

import * as events from 'aws-cdk-lib/aws-events';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { NodejsFunction } from 'aws-cdk-lib/aws-lambda-nodejs';
import { Construct } from 'constructs';
import * as path from 'path';

interface ListenerProps {
    functionName: string;                   // Function name
    lambdaCodePath: string;                 // Path to the Lambda code
    handler: string;                        // Lambda handler function (e.g., "index.handler")
    runtime?: lambda.Runtime;               // Lambda runtime, default to Node.js 16.x
    eventRule: events.Rule;                 // Existing EventBridge rule to add as target
}

export class ListenerLambda extends Construct {
    public readonly lambdaFunction: lambda.Function;

    constructor(scope: Construct, id: string, props: ListenerProps) {
        super(scope, id);

        // Create the Lambda function
        this.lambdaFunction = new NodejsFunction(this, 'ListenerLambda', {
            functionName: props.functionName,
            entry: path.join(props.lambdaCodePath, 'src', 'index.ts'),
            handler: props.handler,
            depsLockFilePath: path.join(props.lambdaCodePath, 'pnpm-lock.yaml'),
            architecture: lambda.Architecture.ARM_64,
            runtime: props.runtime ?? lambda.Runtime.NODEJS_20_X,
            bundling: {
                externalModules: ['aws-sdk'],
                minify: false,
            },
        });

        // Add the Lambda function as a target to the existing EventBridge rule
        props.eventRule.addTarget(new targets.LambdaFunction(this.lambdaFunction));
    }
}

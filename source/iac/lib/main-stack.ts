import * as cdk from 'aws-cdk-lib';
import * as events from 'aws-cdk-lib/aws-events';
import { Construct } from 'constructs';
import { ListenerLambda } from './components/lambda-listener';


export class MainStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const resPrefix = "EventsGate";

    // Define a new EventBridge rule to capture all events
    const allEventsRule = new events.Rule(this, 'AllEventsRule', {
      ruleName: `${resPrefix}DefaultBusRule`,
      eventPattern: {
        account: [this.account]     // Should capture all the events
      },
    });

    new ListenerLambda(this, 'ListenerLambda', {
      functionName: `${resPrefix}ListenerLambda`,
      handler: 'index.handler',
      eventRule: allEventsRule,
      lambdaCodePath: '../backend/listener-lambda/'
    });
  }
}

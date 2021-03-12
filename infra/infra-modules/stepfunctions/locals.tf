locals {
  definition_template = <<EOF
{
  "Comment": "An example of the Amazon States Language for running jobs on Amazon EMR",
  "StartAt": "Create an EMR cluster",
  "States": {
    "Create an EMR cluster": {
      "Type": "Task",
      "Resource": "arn:aws:states:::elasticmapreduce:createCluster.sync",
      "Parameters": {
        "Name": "EMRSparkCluster",
        "VisibleToAllUsers": true,
        "ReleaseLabel": "emr-6.2.0",
        "Applications": [
		  {
            "Name": "Hadoop"
          },
          {
            "Name": "Spark"
          }
        ],
        "ServiceRole": "EMR_DefaultRole",
        "JobFlowRole": "EMR_EC2_DefaultRole",
        "LogUri": "s3://step-functions-emr-${data.aws_caller_identity.current.account_id}/logs/",
        "Instances": {
		  "Ec2SubnetId": "subnet-0452c6fa90d83fcca",	
          "KeepJobFlowAliveWhenNoSteps": true,
          "InstanceFleets": [
            {
              "Name": "EMRSparkMasterFleet",
              "InstanceFleetType": "MASTER",
              "TargetOnDemandCapacity": 1,
              "InstanceTypeConfigs": [
                {
                  "InstanceType": "m5.xlarge"
                }
              ]
            },
            {
              "Name": "EMRSparkCoreFleet",
              "InstanceFleetType": "CORE",
              "TargetOnDemandCapacity": 1,
              "InstanceTypeConfigs": [
                {
                  "InstanceType": "m5.xlarge"
                }
              ]
            }
          ]
        }
      },
      "ResultPath": "$.cluster",
      "Next": "Run first step"
    },
    "Run first step": {
      "Type": "Task",
      "Resource": "arn:aws:states:::elasticmapreduce:addStep.sync",
      "Parameters": {
        "ClusterId.$": "$.cluster.ClusterId",
        "Step": {
          "Name": "My first EMR step",
          "ActionOnFailure": "TERMINATE_CLUSTER",
          "HadoopJarStep": {
            "Jar": "command-runner.jar",
			"Args": ["spark-submit", "--deploy-mode", "cluster", "--class", "org.ussoccer.analytics.dw.DMAData", "s3://step-functions-emr-${data.aws_caller_identity.current.account_id}/libs/spark/dw-assembly-2.0.0.jar"]
          }
        }
      },
	  "Catch": [
		{
		  "ErrorEquals": ["States.ALL"],
		  "Next": "Terminate Cluster"
		}
	  ],
      "ResultPath": "$.firstStep",
      "Next": "Terminate Cluster"
    },
    "Terminate Cluster": {
      "Type": "Task",
      "Resource": "arn:aws:states:::elasticmapreduce:terminateCluster",
      "Parameters": {
        "ClusterId.$": "$.cluster.ClusterId"
      },
      "End": true
    }
  }
}
EOF
}

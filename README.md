# goployer
`goployer` is an application you can use for EC2 deployment. You can deploy in a blue/green mode. goployer only
changes the autoscaling group so that you don't need to create another load balancer or manually attach autoscaling group to target group.
<br><br>

## Demo
![goployer-demo](static/base.gif)

## # Requirements
* You have to create a load balancer and target groups of it which goployer attach a new autoscaling group to. 
* If you want to setup loadbalancer and target group with terraform, then please check this [devopsart workshop](https://devops-art-factory.gitbook.io/devops-workshop/terraform/terraform-resource/computing/elb-+-ec2).
* Please understand how goployer really deploys application before applying to the real environment.
<br>

## # How goployer works
* Here's the steps that goployer executes for deployment
1. Generate new version for current deployment.<br>
If other autoscaling groups of sample application already existed, for example `hello-v001`, then next version will be `hello-v002`
2. Create a new launch template. 
3. Create autoscaling group with launch template from the previous step. A newly created autoscaling group will be automatically attached to the target groups you specified in manifest.
4. Check all instances of all stacks are healty. Until all of them pass healthchecking, it won't go to the next step.
5. (optional) If you add `autoscaling` in manifest, goployer creates autoscaling policies and put these to the autoscaling group. If you use `alarms` with autoscaling, then goployer will also create a cloudwatch alarm for autoscaling policy.
6. After all stacks are deployed, then goployer tries to delete previous versions of the same application.
   Launch templates of previous autoscaling groups are also going to be deleted.
   
<br>

## # Spot Instance
* You can use `spot instance` option with goployer.
* There are two possible ways to use `spot instance`.


`instance_market_options` : You can set spot instance options and with this, you will only use spot instances.
```yaml
    instance_market_options:
      market_type: spot
      spot_options:
        block_duration_minutes: 180
        instance_interruption_behavior: terminate # terminate / stop / hibernate
        max_price: 0.2
        spot_instance_type: one-time # one-time or persistent
```
<br>  
  
`mixed_instances_policy` : You can mix `on-demand` and `spot` together with this setting. 
  
```yaml
    mixed_instances_policy:
      enabled: true
      override_instance_types:
        - c5.large
        - c5.xlarge
      on_demand_percentage: 20
      spot_allocation_strategy: lowest-price
      spot_instance_pools: 3
      spot_max_price: 0.3
```
 
You can see the detailed information in [manifest format](https://goployer.dev/docs/references/manifest/) page.

<br>

## Examples
* You can find few examples of manifest file so that you can test it with this.
```bash
cd examples/manifests
```

## EBS Volume Configuration

Goployer supports advanced EBS volume management with the following features. You can find a complete example in [api-test-example.yaml](examples/manifests/api-test-example.yaml).

### Basic Configuration
```yaml
block_devices:
  - device_name: /dev/xvda
    volume_size: 30
    volume_type: gp3
```

### Delete on Termination
Control whether EBS volumes should be deleted when the instance terminates:
```yaml
block_devices:
  - device_name: /dev/xvda
    volume_size: 30
    volume_type: gp3
    delete_on_termination: false  # optional, defaults to false
```

### Snapshot Support
Create volumes from existing snapshots:
```yaml
block_devices:
  - device_name: /dev/sdf
    snapshot_id: snap-1234567890abcdef0  # optional, for volume creation from snapshot
    volume_size: 100
    volume_type: gp3
    delete_on_termination: true
```

### Encrypted Volumes
Create encrypted volumes with KMS:
```yaml
block_devices:
  - device_name: /dev/sdg
    volume_size: 50
    volume_type: gp3
    encrypted: true
    kms_alias: alias/my-kms-key
    delete_on_termination: false
```

### Configuration Options
- `device_name`: The device name to expose to the instance
- `volume_size`: Size of the volume in GiB
- `volume_type`: Type of EBS volume (gp2, gp3, io1, io2, st1, sc1)
- `delete_on_termination`: Whether to delete the volume on instance termination (default: false)
- `snapshot_id`: ID of the snapshot to create the volume from (optional)
- `encrypted`: Whether to encrypt the volume (default: false)
- `kms_alias`: KMS key alias for encryption (required if encrypted is true)

For a complete example showing how to use these features together, see [api-test-example.yaml](examples/manifests/api-test-example.yaml).



package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"

	"github.com/DevopsArtFactory/goployer/pkg/schemas"
)

func TestMakeLaunchTemplateBlockDeviceMappings(t *testing.T) {
	tests := []struct {
		name     string
		blocks   []schemas.BlockDevice
		expected []*ec2.LaunchTemplateBlockDeviceMappingRequest
	}{
		{
			name: "Basic EBS volume",
			blocks: []schemas.BlockDevice{
				{
					DeviceName: "/dev/xvda",
					VolumeSize: 30,
					VolumeType: "gp2",
				},
			},
			expected: []*ec2.LaunchTemplateBlockDeviceMappingRequest{
				{
					DeviceName: stringPtr("/dev/xvda"),
					Ebs: &ec2.LaunchTemplateEbsBlockDeviceRequest{
						VolumeSize:          int64Ptr(30),
						VolumeType:          stringPtr("gp2"),
						DeleteOnTermination: boolPtr(false),
					},
				},
			},
		},
		{
			name: "EBS volume with DeleteOnTermination",
			blocks: []schemas.BlockDevice{
				{
					DeviceName:          "/dev/xvda",
					VolumeSize:          30,
					VolumeType:          "gp2",
					DeleteOnTermination: true,
				},
			},
			expected: []*ec2.LaunchTemplateBlockDeviceMappingRequest{
				{
					DeviceName: stringPtr("/dev/xvda"),
					Ebs: &ec2.LaunchTemplateEbsBlockDeviceRequest{
						VolumeSize:          int64Ptr(30),
						VolumeType:          stringPtr("gp2"),
						DeleteOnTermination: boolPtr(true),
					},
				},
			},
		},
		{
			name: "EBS volume with valid snapshot",
			blocks: []schemas.BlockDevice{
				{
					DeviceName: "/dev/xvda",
					VolumeSize: 30,
					VolumeType: "gp2",
					SnapshotID: "snap-12345678",
				},
			},
			expected: []*ec2.LaunchTemplateBlockDeviceMappingRequest{
				{
					DeviceName: stringPtr("/dev/xvda"),
					Ebs: &ec2.LaunchTemplateEbsBlockDeviceRequest{
						VolumeSize:          int64Ptr(30),
						VolumeType:          stringPtr("gp2"),
						SnapshotId:          stringPtr("snap-12345678"),
						DeleteOnTermination: boolPtr(false),
					},
				},
			},
		},
		{
			name: "EBS volume with invalid snapshot",
			blocks: []schemas.BlockDevice{
				{
					DeviceName: "/dev/xvda",
					VolumeSize: 30,
					VolumeType: "gp2",
					SnapshotID: "snap-invalid",
				},
			},
			expected: []*ec2.LaunchTemplateBlockDeviceMappingRequest{},
		},
		{
			name: "EBS volume with long snapshot",
			blocks: []schemas.BlockDevice{
				{
					DeviceName: "/dev/xvda",
					VolumeSize: 30,
					VolumeType: "gp2",
					SnapshotID: "snap-1234567890abcdef0",
				},
			},
			expected: []*ec2.LaunchTemplateBlockDeviceMappingRequest{
				{
					DeviceName: stringPtr("/dev/xvda"),
					Ebs: &ec2.LaunchTemplateEbsBlockDeviceRequest{
						VolumeSize:          int64Ptr(30),
						VolumeType:          stringPtr("gp2"),
						SnapshotId:          stringPtr("snap-1234567890abcdef0"),
						DeleteOnTermination: boolPtr(false),
					},
				},
			},
		},
		{
			name: "IOPS volume",
			blocks: []schemas.BlockDevice{
				{
					DeviceName: "/dev/xvda",
					VolumeSize: 30,
					VolumeType: "io1",
					Iops:       3000,
				},
			},
			expected: []*ec2.LaunchTemplateBlockDeviceMappingRequest{
				{
					DeviceName: stringPtr("/dev/xvda"),
					Ebs: &ec2.LaunchTemplateEbsBlockDeviceRequest{
						VolumeSize:          int64Ptr(30),
						VolumeType:          stringPtr("io1"),
						Iops:                int64Ptr(3000),
						DeleteOnTermination: boolPtr(false),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := EC2Client{}
			result := client.MakeLaunchTemplateBlockDeviceMappings(tt.blocks)

			assert.Equal(t, len(tt.expected), len(result))
			for i, expected := range tt.expected {
				assert.Equal(t, *expected.DeviceName, *result[i].DeviceName)
				assert.Equal(t, *expected.Ebs.VolumeSize, *result[i].Ebs.VolumeSize)
				assert.Equal(t, *expected.Ebs.VolumeType, *result[i].Ebs.VolumeType)
				assert.Equal(t, *expected.Ebs.DeleteOnTermination, *result[i].Ebs.DeleteOnTermination)

				if expected.Ebs.SnapshotId != nil {
					assert.Equal(t, *expected.Ebs.SnapshotId, *result[i].Ebs.SnapshotId)
				}
				if expected.Ebs.Iops != nil {
					assert.Equal(t, *expected.Ebs.Iops, *result[i].Ebs.Iops)
				}
			}
		})
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

package controllers

import (
	"testing"

	"github.com/go-logr/logr"
	pvcv1alpha1 "github.com/konveyor/volume-snapshot-mover/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestDataMoverBackupReconciler_CreateResticSecret(t *testing.T) {
	tests := []struct {
		name             string
		dmb              *pvcv1alpha1.DataMoverBackup
		secret, rpsecret *corev1.Secret
		pvc              *corev1.PersistentVolumeClaim
		want             bool
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			name: "Given invalid pvc -> error in restic secter creation",
			dmb: &pvcv1alpha1.DataMoverBackup{
				ObjectMeta: v1.ObjectMeta{
					Name:      "sample-dmb",
					Namespace: "bar",
				},
				Spec: pvcv1alpha1.DataMoverBackupSpec{
					VolumeSnapshotContent: corev1.ObjectReference{
						Name: "sample-snapshot",
					},
					ProtectedNamespace: "foo",
				},
			},
			pvc: &corev1.PersistentVolumeClaim{
				ObjectMeta: v1.ObjectMeta{
					Name:      "sample-snapshot",
					Namespace: "foo",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("10Gi"),
						},
					},
				},
			},
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Name:      resticSecret,
					Namespace: namespace,
				},
				Data: secretData,
			},
			rpsecret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Name:      "sample-dmb-secret",
					Namespace: namespace,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Given valid dmb,restic secret -> successful creation of pvc specific restic secret",
			dmb: &pvcv1alpha1.DataMoverBackup{
				ObjectMeta: v1.ObjectMeta{
					Name:      "sample-dmb",
					Namespace: "bar",
				},
				Spec: pvcv1alpha1.DataMoverBackupSpec{
					VolumeSnapshotContent: corev1.ObjectReference{
						Name: "sample-snapshot",
					},
					ProtectedNamespace: "foo",
				},
			},
			pvc: &corev1.PersistentVolumeClaim{
				ObjectMeta: v1.ObjectMeta{
					Name:      "sample-snapshot-pvc",
					Namespace: "foo",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("10Gi"),
						},
					},
				},
			},
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Name:      resticSecret,
					Namespace: namespace,
				},
				Data: secretData,
			},
			rpsecret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Name:      "sample-dmb-secret",
					Namespace: namespace,
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient, err := getFakeClientFromObjects(tt.dmb, tt.secret, tt.pvc, tt.rpsecret)
			if err != nil {
				t.Errorf("error creating fake client, likely programmer error")
			}
			r := &DataMoverBackupReconciler{
				Client:  fakeClient,
				Scheme:  fakeClient.Scheme(),
				Log:     logr.Discard(),
				Context: newContextForTest(tt.name),
				NamespacedName: types.NamespacedName{
					Namespace: tt.dmb.Spec.ProtectedNamespace,
					Name:      tt.dmb.Name,
				},
				EventRecorder: record.NewFakeRecorder(10),
				req: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: tt.dmb.Namespace,
						Name:      tt.dmb.Name,
					},
				},
			}
			got, err := r.CreateResticSecret(r.Log)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataMoverBackupReconciler.CreateResticSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DataMoverBackupReconciler.CreateResticSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
